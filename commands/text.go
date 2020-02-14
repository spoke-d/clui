package commands

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/task/group"
)

var (
	// ErrShowHelp is a sentinel error message for showing help.
	ErrShowHelp = errors.New("show help")
)

// Text defines a simple text based command, that can be useful for
// creating commands that require more text based explanations.
type Text struct {
	helpText     string
	synopsisText string
	flagSet      *flagset.FlagSet
}

// NewText creates a Command with sane defaults
func NewText(help, synopsis string) *Text {
	return &Text{
		helpText:     help,
		synopsisText: strings.TrimSpace(synopsis),
		flagSet:      flagset.New("text-command", flag.ExitOnError),
	}
}

// FlagSet returns the FlagSet associated with the command. All the flags are
// parsed before running the command.
func (c *Text) FlagSet() *flagset.FlagSet {
	return c.flagSet
}

// Usages returns various usages that can be used for the command.
func (c *Text) Usages() []string {
	return make([]string, 0)
}

// Help should return a long-form help text that includes the command-line
// usage. A brief few sentences explaining the function of the command, and
// the complete list of flags the command accepts.
func (c *Text) Help() string {
	return fmt.Sprintf(`
%q is a placeholder command.

You can see the other sub commands available below.
`, c.helpText)
}

// Synopsis should return a one-line, short synopsis of the command.
// This should be short (50 characters of less ideally).
func (c *Text) Synopsis() string {
	return c.synopsisText
}

// Init is called with all the args required to run a command.
// This is separated from Run, to allow the preperation of a command, before
// it's run.
func (c *Text) Init([]string, bool) error {
	return nil
}

// Run subscribes to the group for executing the various run commands.
// The subscriptions to the group are handled by the callee.
func (c *Text) Run(group *group.Group) {
	group.Add(func() error {
		return ErrShowHelp
	}, Disguard)
}
