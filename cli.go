package clui

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spoke-d/clui/autocomplete"
	"github.com/spoke-d/clui/autocomplete/fsys"
	"github.com/spoke-d/clui/autocomplete/install"
	"github.com/spoke-d/clui/commands"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/group"
	"github.com/spoke-d/clui/help"
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

	// Usages returns various usages that can be used for the command.
	Usages() []string

	// Help should return a long-form help text that includes the command-line
	// usage. A brief few sentences explaining the function of the command, and
	// the complete list of flags the command accepts.
	Help() string

	// Synopsis should return a one-line, short synopsis of the command.
	// This should be short (50 characters of less ideally).
	Synopsis() string

	// Init is called with all the args required to run a command.
	// This is separated from Run, to allow the preperation of a command, before
	// it's run.
	Init([]string, bool) error

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

	args *GlobalArgs
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
		args:          NewGlobalArgs(store),
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

	if err := c.args.Process(args); err != nil {
		return EPerm, err
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
	if c.args.Version() && c.version != "" {
		return c.writeVersion(c.version)
	}

	// Just print the help when only '-h' or '--help' is passed
	if sc := c.args.SubCommand(); c.args.Help() && sc == "" {
		return c.writeHelp(sc)
	}

	// Autocomplete requires the "Name" to be set so that we know what command
	// to setup the autocomplete on.
	if c.name == "" {
		return EPerm, fmt.Errorf("name not set %q", c.name)
	}

	if c.args.RequiresInstall() {
		if err := c.autoCompleter.Install(c.name); err != nil {
			return EPerm, err
		}
	}
	if c.args.RequiresUninstall() {
		if err := c.autoCompleter.Uninstall(c.name); err != nil {
			return EPerm, err
		}
	}

	// Attempt to get the factory function for creating the command
	// implementation. If the command is invalid or blank, it is an error.
	command, ok := c.commands.Get(c.args.SubCommand())
	if !ok {
		return c.writeHelp(c.subCommandParent())
	}

	// Run the command
	if err := command.FlagSet().Parse(c.args.SubCommandArgs()); err != nil {
		return c.commandHelp(command, err.Error())
	}

	// If we've been instructed to just print the help, then print help
	if c.args.Help() {
		return c.commandHelp(command, "")
	}

	// If there is an invalid flag, then error
	if len(c.commandFlags) > 0 {
		return c.commandHelp(command, "")
	}

	// Remove the flags, those are handled by the flagset.
	var subCmdArgs []string
	for _, v := range c.args.SubCommandArgs() {
		if !strings.HasPrefix(v, "-") {
			subCmdArgs = append(subCmdArgs, v)
		}
	}
	if err := command.Init(subCmdArgs, c.args.Debug()); err != nil {
		return c.commandHelp(command, err.Error())
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
		return c.commandHelp(command, "")
	case nil:
		return EOK, nil
	default:
		return c.commandHelp(command, err.Error())
	}
}

// subCommandParent returns the parent of this subCommand, if there is one.
// Returns empty string ("") if this isn't a parent.
func (c *CLI) subCommandParent() string {
	// get the subCommand, if it is "", just return
	sub := c.args.SubCommand()
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

func (c *CLI) writeHelp(command string) (Errno, error) {
	children, err := FindChildren(c.commands, command)
	if err != nil {
		return EPerm, errors.WithStack(err)
	}

	shims := make(map[string]help.Command, len(children))
	for k, v := range children {
		shims[k] = v
	}

	subCommand := c.args.SubCommand()

	var hint string
	if close, ok := c.commands.GetClosestName(subCommand); ok {
		hint = close
	}

	var header string
	if subCommand == "" {
		header = c.header
	}

	fn := help.BasicFunc(fmt.Sprintf("%s %s", c.name, subCommand))
	res, err := fn(
		help.OptionCommands(shims),
		help.OptionHeader(header),
		help.OptionHint(hint),
		help.OptionColor(!c.args.RequiresNoColor()),
		help.OptionTemplate(help.BasicHelpTemplate),
	)
	if err != nil {
		return EPerm, errors.WithStack(err)
	}
	c.ui.Info(res)
	return EOK, nil
}

func (c *CLI) commandHelp(command Command, operatorErr string) (Errno, error) {
	subCommand := c.args.SubCommand()

	children, err := FindChildren(c.commands, subCommand)
	if err != nil {
		return EPerm, errors.WithStack(err)
	}

	shims := make(map[string]help.Command, len(children))
	for k, v := range children {
		shims[k] = v
	}

	var hint string
	if close, ok := c.commands.GetClosestName(subCommand); ok && close != subCommand {
		hint = close
	}

	var header string
	if subCommand == "" {
		header = c.header
	}

	flags, err := commandFlags(command.FlagSet())
	if err != nil {
		return EPerm, errors.WithStack(err)
	}

	fn := help.BasicFunc(fmt.Sprintf("%s %s", c.name, subCommand))
	res, err := fn(
		help.OptionCommands(shims),
		help.OptionHeader(header),
		help.OptionHint(hint),
		help.OptionColor(!c.args.RequiresNoColor()),
		help.OptionTemplate(help.CommandHelpTemplate),
		help.OptionHelp(command.Help()),
		help.OptionFlags(flags),
		help.OptionUsages(command.Usages()),
		help.OptionErr(operatorErr),
	)
	if err != nil {
		return EPerm, errors.WithStack(err)
	}
	c.ui.Info(res)
	return EOK, nil
}

func (c *CLI) writeVersion(s string) (Errno, error) {
	template := ui.NewTemplate(TemplateVersion)
	return EOK, c.ui.Output(template, struct {
		Version string
	}{
		Version: s,
	})
}

func commandFlags(flags *flagset.FlagSet) ([]string, error) {
	type flagType struct {
		Name     string
		Usage    string
		Defaults string
	}

	template := ui.NewTemplate(TemplateFlags, ui.OptionName("flags"))
	var allFlags []*flag.Flag
	flags.VisitAll(func(f *flag.Flag) {
		allFlags = append(allFlags, f)
	})

	data := make([]string, len(allFlags))
	for k, v := range allFlags {
		res, err := template.Render(flagType{
			Name:     fmt.Sprintf("--%s", v.Name),
			Usage:    v.Usage,
			Defaults: v.DefValue,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		data[k] = strings.TrimSpace(res)
	}

	return data, nil
}
