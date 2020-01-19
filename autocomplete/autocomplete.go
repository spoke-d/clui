package autocomplete

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spoke-d/clui/autocomplete/args"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/radix"
)

const (
	envComplete = "COMP_LINE"
	envDebug    = "COMP_DEBUG"
)

// Installer is an interface to be implemented to perform the autocomplete
// installation and un-installation with a CLI.
//
// This interface is not exported because it only exists for unit tests
// to be able to test that the installation is called properly.
type Installer interface {
	Install(string) error
	Uninstall(string) error
}

// Group describes an abstraction for walking over a group of commands by
// the prefix provided.
type Group interface {
	WalkPrefix(string, radix.WalkFn)
}

// Command represents an abstraction of command.
type Command interface {
	FlagSet() *flagset.FlagSet
}

// AutoCompleteOptions represents a way to set optional values to a autocomplete
// option.
// The AutoCompleteOptions shows what options are available to change.
type AutoCompleteOptions interface {
	SetInstaller(Installer)
	SetGroup(Group)
}

// AutoCompleteOption captures a tweak that can be applied to the AutoComplete.
type AutoCompleteOption func(AutoCompleteOptions)

type autocomplete struct {
	installer Installer
	group     Group
}

func (s *autocomplete) SetInstaller(i Installer) {
	s.installer = i
}

func (s *autocomplete) SetGroup(g Group) {
	s.group = g
}

// OptionInstaller allows the setting a installer option to configure
// the autocomplete.
func OptionInstaller(i Installer) AutoCompleteOption {
	return func(opt AutoCompleteOptions) {
		opt.SetInstaller(i)
	}
}

// OptionGroup allows the setting a installer option to configure
// the autocomplete.
func OptionGroup(g Group) AutoCompleteOption {
	return func(opt AutoCompleteOptions) {
		opt.SetGroup(g)
	}
}

// AutoComplete defines a way to predict and complete arguments passed
// in to the CLI
type AutoComplete struct {
	installer Installer
	group     Group
}

// New creates a new AutoComplete with the correct dependencies.
func New(options ...AutoCompleteOption) *AutoComplete {
	opt := new(autocomplete)
	for _, option := range options {
		option(opt)
	}

	return &AutoComplete{
		installer: opt.installer,
		group:     opt.group,
	}
}

// Install a command into the host using the Installer.
// Returns an error if there is an error whilst installing.
func (a *AutoComplete) Install(cmd string) error {
	return a.installer.Install(cmd)
}

// Uninstall the command from the host using the Installer.
// Returns an error if there is an error whilst it's uninstalling.
func (a *AutoComplete) Uninstall(cmd string) error {
	return a.installer.Uninstall(cmd)
}

// Complete a command from completion line in environment variable,
// and print out the complete options.
// Returns success if the completion ran or if the cli matched
// any of the given flags, false otherwise
func (a *AutoComplete) Complete(line string) ([]string, bool) {
	// If the line is empty, don't even attempt to complete on it.
	if line == "" {
		return nil, false
	}

	// If we've got something attempt to predict the command.
	var (
		matches []string

		args    = args.New(line)
		options = a.Predict(args)
	)

	for _, opt := range options {
		if strings.HasPrefix(opt, args.Last()) {
			matches = append(matches, opt)
		}
	}

	return matches, true
}

// Predict returns all possible predictions for args according to the command.
func (a *AutoComplete) Predict(v *args.Args) []string {
	var (
		options   []string
		potential []pair
		args      = strings.Join(v.AllCommands(), " ")
	)
	a.group.WalkPrefix(args, func(s string, cmd radix.Value) bool {
		if c, ok := cmd.(Command); ok {
			potential = append(potential, pair{
				Name:    s,
				Command: c,
			})
			return true
		}
		return false
	})

	if len(potential) > 0 {
		if isFlag := strings.HasPrefix(v.Last(), "-"); isFlag {
			// Check if the potential cmd is an exact match
			var pairs []pair
			for _, pair := range potential {
				if pair.Name == strings.Join(v.AllCommands(), " ") {
					pairs = append(pairs, pair)
				}
			}

			// find out what those flags are
			for _, pair := range pairs {
				opts, only := predictFlag(pair.Command, v)
				if only {
					return opts
				}
				options = append(options, opts...)
			}
		} else {
			// auto complete the command name
			for _, pair := range potential {
				parts := strings.Split(pair.Name, " ")
				if len(parts) >= 1 {
					options = append(options, parts[len(parts)-1])
				}
			}
		}
	} else {
		// TODO: auto complete files
	}
	return options
}

type pair struct {
	Name    string
	Command Command
}

func predictFlag(cmd Command, a *args.Args) ([]string, bool) {
	flagset := cmd.FlagSet()
	flagName := strings.TrimLeft(strings.TrimSpace(a.Last()), "-")
	if flag := flagset.Lookup(flagName); flag != nil {
		return []string{fmt.Sprintf("--%s", flag.Name)}, true
	}

	var options []string
	flagset.VisitAll(func(f *flag.Flag) {
		options = append(options, fmt.Sprintf("--%s", f.Name))
	})

	return options, false
}
