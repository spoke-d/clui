package commands

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/pkg/errors"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/radix"
	"github.com/spoke-d/task/group"
)

// Runnable allows the shell to interact with a given runnable type.
type Runnable interface {
	// Run runs the actual CLI bases on the arguments given.
	Run(args []string) (int, error)
}

// Store holds the command prefixes to able to walk over.
type Store interface {
	// WalkPrefix is used to walk the tree under a prefix
	WalkPrefix(prefix string, fn radix.WalkFn)
}

// Shell defines a REPL that can be interactively accessed.
type Shell struct {
	flagSet *flagset.FlagSet
	runner  Runnable
	group   Store
}

// NewShell creates a REPL from a runnable and a command store.
func NewShell(runner Runnable, group Store) *Shell {
	return &Shell{
		flagSet: flagset.New("text-command", flag.ContinueOnError),
		runner:  runner,
		group:   group,
	}
}

// FlagSet returns the FlagSet associated with the command. All the flags are
// parsed before running the command.
func (c *Shell) FlagSet() *flagset.FlagSet {
	return c.flagSet
}

// Usages returns various usages that can be used for the command.
func (c *Shell) Usages() []string {
	return make([]string, 0)
}

var firstPrompt = `
Welcome to the interactive shell.
Type "help" to see a list of available commands.
Type ^D or ^C to quit.
`[1:]

// Help should return a long-form help text that includes the command-line
// usage. A brief few sentences explaining the function of the command, and
// the complete list of flags the command accepts.
func (c *Shell) Help() string {
	return `
The shell command allows REPL functionality for interacting with
the commands for the given CLI. All arguments and flags are then
parsed and forwared to the correct command.

Type ^D or ^C to exit the shell.`
}

// Synopsis should return a one-line, short synopsis of the command.
// This should be short (50 characters of less ideally).
func (c *Shell) Synopsis() string {
	return "Shell command invokes a REPL for command interaction."
}

// Init is called with all the args required to run a command.
// This is separated from Run, to allow the preperation of a command, before
// it's run.
func (c *Shell) Init([]string, CommandContext) error {
	return nil
}

// Run subscribes to the group for executing the various run commands.
// The subscriptions to the group are handled by the callee.
func (c *Shell) Run(group *group.Group) {
	group.Add(func(context.Context) error {
		history, err := ioutil.TempFile("", "convoy-shell")
		if err != nil {
			return errors.WithStack(err)
		}
		defer history.Close()

		line, err := readline.NewEx(&readline.Config{
			HistoryFile:     history.Name(),
			Stdin:           readline.NewCancelableStdin(os.Stdin),
			Stdout:          os.Stdout,
			Stderr:          os.Stderr,
			Prompt:          "\033[31m»\033[0m ",
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
			AutoComplete:    generateCompletions(c.group),
		})
		if err != nil {
			return errors.WithStack(err)
		}
		defer line.Close()

		first := true
		for {
			if first {
				fmt.Fprintln(os.Stdout, firstPrompt)
				first = false
			}

			data, err := line.Readline()
			if err == readline.ErrInterrupt {
				if len(data) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}

			if cmd := strings.TrimSpace(data); cmd == "help commands" {
				fmt.Fprintln(os.Stdout, listAllCommands(c.group))
				continue
			} else if cmd == "help" {
				fmt.Fprintln(os.Stdout, "Type ^D or ^C to exit the shell.")
				continue
			} else if cmd == "exit" {
				return nil
			}

			if _, err := c.runner.Run(strings.Split(data, " ")); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		return nil
	}, Disguard)
}

func listAllCommands(group Store) string {
	var commands []string
	group.WalkPrefix("", func(name string, value radix.Value) bool {
		commands = append(commands, name)
		return false
	})
	return strings.Join(commands, "\n")
}

func generateCompletions(group Store) *readline.PrefixCompleter {
	nodes := node{
		name:     "",
		children: make(map[string]node),
	}
	group.WalkPrefix("", func(name string, value radix.Value) bool {
		names := strings.Split(name, " ")
		parent := nodes
		for _, n := range names {
			if _, ok := parent.children[n]; ok {
				parent = parent.children[n]
				continue
			}
			m := node{
				name:     n,
				children: make(map[string]node),
			}

			parent.children[n] = m
			parent = m
		}
		return false
	})
	var recursive func(nodes node) []readline.PrefixCompleterInterface
	recursive = func(nodes node) []readline.PrefixCompleterInterface {
		if len(nodes.children) == 0 {
			return []readline.PrefixCompleterInterface{
				readline.PcItem(nodes.name),
			}
		}

		items := make([]readline.PrefixCompleterInterface, 0)
		for _, node := range nodes.children {
			items = append(items, recursive(node)...)
		}
		return []readline.PrefixCompleterInterface{
			readline.PcItem(nodes.name, items...),
		}
	}

	items := recursive(nodes)
	if len(items) != 1 {
		return readline.NewPrefixCompleter()
	}
	generated := items[0].(*readline.PrefixCompleter).Children
	root := []readline.PrefixCompleterInterface{
		readline.PcItem("help",
			readline.PcItem("commands"),
		),
		readline.PcItem("exit"),
	}
	return readline.NewPrefixCompleter(append(root, generated...)...)
}

type node struct {
	name     string
	children map[string]node
}
