package clui

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spoke-d/clui/group"
)

// GlobalArgs is used to construct the arguments used for the CLI.
// The global arguments are then passed to the command once found, without the
// global flags.
type GlobalArgs struct {
	commands     *group.Group
	commandFlags []string

	subCommand      string
	subCommandArgs  []string
	subCommandFlags []string

	isHelp, isVersion, isDebug, isDevMode bool
	requiresInstall, requiresUninstall    bool
	requiresNoColor                       bool
	requiresNoSubKeys                     bool
}

// NewGlobalArgs creates a new GlobalArgs type for processing arguments passed
// to the cli.
func NewGlobalArgs(commands *group.Group) *GlobalArgs {
	return &GlobalArgs{
		commands: commands,
	}
}

// SubCommand returns the sub command name.
func (a *GlobalArgs) SubCommand() string {
	return a.subCommand
}

// SubCommandArgs returns the sub command arguments.
func (a *GlobalArgs) SubCommandArgs() []string {
	return a.subCommandArgs
}

// SubCommandFlags returns the sub command flags.
func (a *GlobalArgs) SubCommandFlags() []string {
	return a.subCommandFlags
}

// CommandFlags returns the command arguments.
func (a *GlobalArgs) CommandFlags() []string {
	return a.commandFlags
}

// Help returns if the operator has passed the help flag.
func (a *GlobalArgs) Help() bool {
	return a.isHelp
}

// Version returns if the operator has passed the version flag.
func (a *GlobalArgs) Version() bool {
	return a.isVersion
}

// Debug returns if the operator has passed the debug flag.
func (a *GlobalArgs) Debug() bool {
	return a.isDebug
}

// DevMode returns if the operator has passed the devMode flag.
func (a *GlobalArgs) DevMode() bool {
	return a.isDevMode
}

// RequiresInstall returns if the operator has passed the requires install flag.
func (a *GlobalArgs) RequiresInstall() bool {
	return a.requiresInstall
}

// RequiresUninstall returns if the operator has passed the requires uninstall
// flag.
func (a *GlobalArgs) RequiresUninstall() bool {
	return a.requiresUninstall
}

// RequiresNoColor returns if the operator has passed the no color output flag.
func (a *GlobalArgs) RequiresNoColor() bool {
	return a.requiresNoColor
}

// RequiresNoSubKeys returns if the help should include subkeys.
func (a *GlobalArgs) RequiresNoSubKeys() bool {
	return a.requiresNoSubKeys
}

// Process consumes the arguments and correctly separates them between global
// flags and arguments and command flags and arguments.
func (a *GlobalArgs) Process(args []string) error {
	var processed []string
	// first remove the default
	for _, arg := range args {
		if arg == "--" {
			break
		}

		switch arg {
		case "-h", "-help", "--help":
			a.isHelp = true
		case "-v", "-version", "--version":
			a.isVersion = true
		case "--debug":
			a.isDebug = true
		case "--dev-mode":
			a.isDevMode = true
		case "--no-color":
			a.requiresNoColor = true
		case "--no-sub-keys":
			a.requiresNoSubKeys = true
		case "--autocomplete-install":
			a.requiresInstall = true
		case "--autocomplete-uninstall":
			a.requiresUninstall = true
		default:
			processed = append(processed, arg)
		}
	}

	if a.requiresInstall && a.requiresUninstall {
		return errors.Errorf("both autocomplete flags can not be used at the same time")
	}

	for i, arg := range processed {
		if a.subCommand == "" {
			if arg != "" && arg[0] == '-' {
				// Record the arg...
				a.commandFlags = append(a.commandFlags, arg)
			}
		}

		// If we didn't find a subCommand yet and this is the first non-flag
		// argument, then this is our subCommand.
		if a.subCommand == "" && arg != "" && arg[0] != '-' {
			a.subCommand = arg
			if a.commands.Nested() {
				// If the command has a space in it, then it is invalid.
				// Set a blank command so that it fails.
				if strings.ContainsRune(arg, ' ') {
					a.subCommand = ""
					return nil
				}

				// Determine the argument we look to end subCommands.
				// We look at all arguments until one has a space. This
				// disallows commands like: ./cli foo "bar baz".
				// An argument with a space is always an argument.
				var j int
				for k, v := range processed[i:] {
					if strings.ContainsRune(v, ' ') {
						break
					}
					j = i + k + 1
				}

				// Nested CLI the subCommand is actually the entire arg list up
				// to a flag that is still a valid subCommand.
				searchKey := strings.Join(processed[i:j], " ")
				k, ok := a.commands.LongestPrefix(searchKey)
				if ok {
					// k could be a prefix that doesn't contain the full command
					// such as "foo", instead of "foobar", so we need to verify
					// that we have an entire key. To do that, we look for an
					// ending in a space of end of a string.
					verify, err := regexp.Compile(regexp.QuoteMeta(k) + `( |$)`)
					if err != nil {
						return err
					}
					if verify.MatchString(searchKey) {
						a.subCommand = k
						i += strings.Count(k, " ")
					}
				}
			}

			// The remaining processed the subCommand arguments

			a.subCommandArgs = removeFlags(processed[i+1:])
			a.subCommandFlags = removeNonFlags(processed[i+1:])
		}
	}

	return nil
}

func removeFlags(args []string) []string {
	var result []string
	for _, v := range args {
		if v == "-" || !strings.HasPrefix(v, "-") {
			result = append(result, v)
		}
	}
	return result
}

func removeNonFlags(args []string) []string {
	var result []string
	for _, v := range args {
		if strings.HasPrefix(v, "-") {
			result = append(result, v)
		}
	}
	return result
}
