package clui

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spoke-d/clui/autocomplete"
	"github.com/spoke-d/clui/autocomplete/fsys"
	"github.com/spoke-d/clui/autocomplete/install"
	"github.com/spoke-d/clui/commands"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/group"
	"github.com/spoke-d/clui/help"
	"github.com/spoke-d/clui/radix"
	"github.com/spoke-d/clui/ui"
	task "github.com/spoke-d/task/group"
)

// UI is an interface for interacting with the terminal, or "interface"
// of a CLI.
type UI interface {
	// Ask asks the user for input using the given query. The response is
	// returned as the given string, or an error.
	Ask(string) (string, error)

	// AskSecret asks the user for input using the given query, but does not echo
	// the keystrokes to the terminal.
	AskSecret(string) (string, error)

	// Output is called for normal standard output.
	Output(*ui.Template, interface{}) error

	// Info is called for information related to the previous output.
	// In general this may be the exact same as Output, but this gives
	// UI implementors some flexibility with output formats.
	Info(string)

	// Error is used for any error messages that might appear on standard
	// error.
	Error(string)
}

// Command is a runnable sub-command of CLI.
type Command interface {

	// Flags returns the FlagSet associated with the command. All the flags are
	// parsed before running the command.
	FlagSet() *flagset.FlagSet

	// Help should return a long-form help text that includes the command-line
	// usage. A brief few sentences explaining the function of the command, and
	// the complete list of flags the command accepts.
	Help() string

	// Synopsis should return a one-line, short synopsis of the command.
	// This should be short (50 characters of less ideally).
	Synopsis() string

	// Run should run the actual command with the given CLI instance and
	// command-line arguments. It should return the exit status when it is
	// finished.
	//
	// There are a handful of special exit codes that can return documented
	// behavioral changes.
	Run(*task.Group)
}

// AutoCompleter is an interface to be implemented to perform the autocomplete
// installation and un-installation with a CLI.
//
// This interface is not exported because it only exists for unit tests
// to be able to test that the installation is called properly.
type AutoCompleter interface {
	// Complete a command from completion line in environment variable,
	// and print out the complete options.
	// Returns success if the completion ran or if the cli matched
	// any of the given flags, false otherwise
	Complete(string) ([]string, bool)

	// Install a command into the host using the Installer.
	// Returns an error if there is an error whilst installing.
	Install(string) error

	// Uninstall the command from the host using the Installer.
	// Returns an error if there is an error whilst it's uninstalling.
	Uninstall(string) error
}

// CLIOptions represents a way to set optional values to a CLI option.
// The CLIOptions shows what options are available to change.
type CLIOptions interface {
	SetHelpFunc(help.Func)
	SetAutoCompleter(AutoCompleter)
	SetUI(UI)
	SetFileSystem(fsys.FileSystem)
}

// CLIOption captures a tweak that can be applied to the CLI.
type CLIOption func(CLIOptions)

type cli struct {
	helpFunc      help.Func
	autoCompleter AutoCompleter
	fileSystem    fsys.FileSystem
	ui            UI
}

func (s *cli) SetHelpFunc(p help.Func) {
	s.helpFunc = p
}

func (s *cli) HelpFunc(name string) help.Func {
	if s.helpFunc == nil {
		return help.BasicFunc(name)
	}
	return s.helpFunc
}

func (s *cli) SetAutoCompleter(p AutoCompleter) {
	s.autoCompleter = p
}

func (s *cli) SetFileSystem(p fsys.FileSystem) {
	s.fileSystem = p
}

func (s *cli) AutoCompleter(group *group.Group, fs fsys.FileSystem) AutoCompleter {
	if s.autoCompleter == nil {
		user, err := install.CurrentUser()
		if err != nil {
			return nil
		}
		installer, err := install.New(
			install.OptionShell(install.Bash(install.OptionUser(user), install.OptionFileSystem(fs))),
			install.OptionShell(install.Zsh(install.OptionUser(user), install.OptionFileSystem(fs))),
		)
		if err != nil {
			return nil
		}
		return autocomplete.New(autocomplete.OptionGroup(group),
			autocomplete.OptionInstaller(installer),
		)
	}
	return s.autoCompleter
}

func (s *cli) SetUI(p UI) {
	s.ui = p
}

func (s *cli) UI() UI {
	if s.ui == nil {
		return ui.NewBasicUI(os.Stdin, os.Stdout, os.Stderr)
	}
	return s.ui
}

// OptionHelpFunc allows the setting a HelpFunc option to configure the cli.
func OptionHelpFunc(i help.Func) CLIOption {
	return func(opt CLIOptions) {
		opt.SetHelpFunc(i)
	}
}

// OptionAutoCompleter allows the setting a AutoCompleter option to configure
// the cli.
func OptionAutoCompleter(i AutoCompleter) CLIOption {
	return func(opt CLIOptions) {
		opt.SetAutoCompleter(i)
	}
}

// OptionUI allows the setting a UI option to configure the cli.
func OptionUI(i UI) CLIOption {
	return func(opt CLIOptions) {
		opt.SetUI(i)
	}
}

// OptionFileSystem allows the setting a FileSystem option to configure the cli.
func OptionFileSystem(i fsys.FileSystem) CLIOption {
	return func(opt CLIOptions) {
		opt.SetFileSystem(i)
	}
}

// CommandFn defines a function for constructing a command.
type CommandFn func(UI) Command

// CLI contains the state necessary to run commands and parse the command line
// arguments
//
// CLI also supports nested subCommands, such as "cli foo bar". To use nested
// subCommands, the key in the Commands mapping below contains the full
// subCommand.
// In this example, it would be "foo bar"
type CLI struct {
	name    string
	version string
	header  string

	ui            UI
	autoCompleter AutoCompleter

	// HelpFunc and HelpWriter are used to output help information, if
	// requested.
	//
	// HelpFunc is the function called to generate the generic help text that is
	// shown. If help must be shown for the CLI that doesn't pertain to a
	// specific command.
	//
	// HelpWriter is the Writer where the help text is outputted to. If not
	// specified, it will default to Stderr.
	helpFunc help.Func

	commands     *group.Group
	commandFlags []string

	subCommand     string
	subCommandArgs []string

	// These are special global flags
	isHelp, isVersion, isDebug         bool
	requiresInstall, requiresUninstall bool
	requiresNoColor                    bool
}

// New returns a new CLI instance with sensible default.
func New(name, version, header string, options ...CLIOption) *CLI {
	opt := new(cli)
	for _, option := range options {
		option(opt)
	}

	store := group.New(group.OptionPlaceHolder(func(s string) group.Command {
		return commands.NewText(s, TemplatePlaceHolder)
	}))

	return &CLI{
		name:          name,
		version:       version,
		header:        header,
		ui:            opt.UI(),
		helpFunc:      opt.HelpFunc(name),
		commands:      store,
		autoCompleter: opt.AutoCompleter(store, opt.fileSystem),
	}
}

// Add inserts a new command to the CLI.
func (c *CLI) Add(key string, cmdFn CommandFn) error {
	return c.commands.Add(key, cmdFn(c.ui))
}

// Run runs the actual CLI bases on the arguments given.
func (c *CLI) Run(args []string) (Errno, error) {
	if err := c.commands.Process(); err != nil {
		return EPerm, err
	}

	if err := c.processArgs(args); err != nil {
		return EPerm, err
	}

	if c.requiresInstall && c.requiresUninstall {
		return EPerm, fmt.Errorf("both autocomplete flags can not be used at the same time")
	}

	// If this is a autocompletion request, satisfy it. This must be called
	// first before anything else since its possible to be autocompleting
	// -help or -version or other flags and we want to show completions
	// and not actually write the help or version.
	// TODO: Get this from options
	if commands, ok := c.autoCompleter.Complete(autocomplete.TerminalLine()); ok {
		template := ui.NewTemplate(TemplateComplete)
		return EOK, c.ui.Output(template, commands)
	}

	// Just show the version and exit if instructed.
	if c.isVersion && c.version != "" {
		return c.writeVersion(c.version)
	}

	// Just print the help when only '-h' or '--help' is passed
	if sc := c.subCommand; c.isHelp && sc == "" {
		return c.writeHelp(sc, helpErrCode{
			Success: 0,
			Failure: EPerm,
		})
	}

	// Autocomplete requires the "Name" to be set so that we know what command
	// to setup the autocomplete on.
	if c.name == "" {
		return EPerm, fmt.Errorf("name not set %q", c.name)
	}

	if c.requiresInstall {
		if err := c.autoCompleter.Install(c.name); err != nil {
			return EPerm, err
		}
	}
	if c.requiresUninstall {
		if err := c.autoCompleter.Uninstall(c.name); err != nil {
			return EPerm, err
		}
	}

	// Attempt to get the factory function for creating the command
	// implementation. If the command is invalid or blank, it is an error.
	command, ok := c.commands.Get(c.subCommand)
	if !ok {
		return c.writeHelp(c.subCommandParent(), helpErrCode{
			Success: EKeyExpired,
			Failure: EPerm,
		})
	}

	// Run the command
	if err := command.FlagSet().Parse(c.subCommandArgs); err != nil {
		return EPerm, c.commandHelp(command, err.Error())
	}

	// If we've been instructed to just print the help, then print help
	if c.isHelp {
		return EOK, c.commandHelp(command, "")
	}

	// If there is an invalid flag, then error
	if len(c.commandFlags) > 0 {
		return EPerm, c.commandHelp(command, "")
	}

	// Create a new group context to run.
	g := task.NewGroup()

	// Subscribe all the actions to the group.
	task.Block(g)
	command.Run(g)
	task.Interrupt(g)

	// Run the group
	switch err := g.Run(); err {
	case commands.ErrShowHelp:
		return EPerm, c.commandHelp(command, "")
	case nil:
		return EOK, nil
	default:
		return EPerm, c.commandHelp(command, err.Error())
	}
}

// IsVersion returns whether or not the version flag is present within the
// arguments.
func (c *CLI) IsVersion() bool {
	return c.isVersion
}

// IsHelp returns whether or not the help flag is present within the arguments.
func (c *CLI) IsHelp() bool {
	return c.isHelp
}

func (c *CLI) processArgs(args []string) error {
	var processed []string
	// first remove the default
	for _, arg := range args {
		if arg == "--" {
			break
		}

		switch arg {
		case "-h", "-help", "--help":
			c.isHelp = true
		case "-v", "-version", "--version":
			c.isVersion = true
		case "--debug":
			c.isDebug = true
		case "--no-color":
			c.requiresNoColor = true
		case "--autocomplete-install":
			c.requiresInstall = true
		case "--autocomplete-uninstall":
			c.requiresUninstall = true
		default:
			processed = append(processed, arg)
		}
	}

	for i, arg := range processed {
		if c.subCommand == "" {
			if arg != "" && arg[0] == '-' {
				// Record the arg...
				c.commandFlags = append(c.commandFlags, arg)
			}
		}

		// If we didn't find a subCommand yet and this is the first non-flag
		// argument, then this is our subCommand.
		if c.subCommand == "" && arg != "" && arg[0] != '-' {
			c.subCommand = arg
			if c.commands.Nested() {
				// If the command has a space in it, then it is invalid.
				// Set a blank command so that it fails.
				if strings.ContainsRune(arg, ' ') {
					c.subCommand = ""
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
				k, ok := c.commands.LongestPrefix(searchKey)
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
						c.subCommand = k
						i += strings.Count(k, " ")
					}
				}
			}

			// The remaining processed the subCommand arguments
			c.subCommandArgs = processed[i+1:]
		}
	}

	// If we never found a subCommand and support a default command, then
	// switch to using that
	if c.subCommand == "" {
		if _, ok := c.commands.Get(""); ok {
			x := c.commandFlags
			x = append(x, c.subCommandArgs...)
			c.commandFlags = nil
			c.subCommandArgs = x
		}
	}

	return nil
}

// helpCommands returns the subCommands for the HelpFunc argument.
// This will only contain immediate subCommands.
func (c *CLI) helpCommands(prefix string) (map[string]Command, error) {
	// if our prefix isn't empty, make sure it ends in ' '
	if prefix != "" && prefix[len(prefix)-1] != ' ' {
		prefix += " "
	}

	// Get all the subkeys of this command
	var keys []string
	c.commands.WalkPrefix(prefix, func(k string, v radix.Value) bool {
		// Ignore any sub-sub keys, i.e. "foo bar baz" when we want "foo bar"
		if !strings.Contains(k[len(prefix):], " ") {
			keys = append(keys, k)
		}

		return false
	})

	// For each of the keys return that in the map
	res := make(map[string]Command, len(keys))
	for _, k := range keys {
		cmd, ok := c.commands.Get(k)
		if !ok {
			return nil, fmt.Errorf("not found: %q", k)
		}
		res[k] = cmd
	}

	return res, nil
}

func (c *CLI) commandHelp(command Command, err string) error {
	var buf bytes.Buffer
	writer := tabwriter.NewWriter(&buf, 24, 4, 4, ' ', 0)

	var hint string
	if close, ok := c.commands.GetClosestName(c.subCommand); ok && close != c.subCommand {
		hint = close
	}

	var flags []string
	command.FlagSet().VisitAll(func(f *flag.Flag) {
		flags = append(flags, fmt.Sprintf("--%s	%s (default %q)", f.Name, f.Usage, f.DefValue))
	})

	template := ui.NewTemplate(TemplateHelp,
		ui.OptionName("help"),
		ui.OptionColor(!c.requiresNoColor),
	)
	if err := template.Write(writer, struct {
		Err   string
		Help  string
		Name  string
		Hint  string
		Flags []string
	}{
		Err:   err,
		Help:  command.Help(),
		Name:  c.subCommand,
		Hint:  hint,
		Flags: flags,
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := writer.Flush(); err != nil {
		return errors.WithStack(err)
	}

	c.ui.Info(strings.TrimSpace(buf.String()) + "\n")

	return nil
}

// subCommandParent returns the parent of this subCommand, if there is one.
// Returns empty string ("") if this isn't a parent.
func (c *CLI) subCommandParent() string {
	// get the subCommand, if it is "", just return
	sub := c.subCommand
	if sub == "" {
		return sub
	}

	// Clear any trailing spaces and find the last space
	sub = strings.TrimRight(sub, " ")
	idx := strings.LastIndex(sub, " ")
	if idx == -1 {
		// No space means our parent is the root
		return ""
	}
	return sub[:idx]
}

type helpErrCode struct {
	Success, Failure Errno
}

func (c *CLI) writeHelp(command string, errCode helpErrCode) (Errno, error) {
	cmds, err := c.helpCommands(command)
	if err != nil {
		return errCode.Failure, err
	}

	shims := make(map[string]help.Command, len(cmds))
	for k, v := range cmds {
		shims[k] = v
	}

	var hint string
	if close, ok := c.commands.GetClosestName(c.subCommand); ok {
		hint = close
	}

	var header string
	if c.subCommand == "" {
		header = c.header
	}

	res, err := c.helpFunc(help.OptionCommands(shims),
		help.OptionHeader(header),
		help.OptionHint(hint),
		help.OptionColor(!c.requiresNoColor),
	)
	if err != nil {
		return errCode.Failure, err
	}
	c.ui.Info(res)
	return errCode.Success, nil
}

func (c *CLI) writeVersion(s string) (Errno, error) {
	template := ui.NewTemplate(TemplateVersion)
	return EOK, c.ui.Output(template, struct {
		Version string
	}{
		Version: s,
	})
}
