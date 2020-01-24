package help

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spoke-d/clui/ui"
)

// Command represents an abstraction of command.
type Command interface {
	Synopsis() string
}

// HelpOptions represents a way to set optional values to a autocomplete
// option.
// The HelpOptions shows what options are available to change.
type HelpOptions interface {
	SetHeader(string)
	SetHint(string)
	SetHelp(string)
	SetErr(string)
	SetCommands(map[string]Command)
	SetFlags([]string)
	SetFormat(string)
	SetColor(bool)
	SetTemplate(string)
}

// HelpOption captures a tweak that can be applied to the Help.
type HelpOption func(HelpOptions)

// HelpOptions defines options for overriding help rendering.
type help struct {
	header   string
	hint     string
	help     string
	err      string
	commands map[string]Command
	flags    []string
	format   string
	color    bool
	template string
}

func (s *help) SetHeader(p string) {
	s.header = p
}

func (s *help) SetHint(p string) {
	s.hint = p
}

func (s *help) SetHelp(p string) {
	s.help = p
}

func (s *help) SetErr(p string) {
	s.err = p
}

func (s *help) SetCommands(p map[string]Command) {
	s.commands = p
}

func (s *help) SetFlags(p []string) {
	s.flags = p
}

func (s *help) SetFormat(p string) {
	s.format = p
}

func (s *help) SetColor(p bool) {
	s.color = p
}

func (s *help) SetTemplate(p string) {
	s.template = p
}

// OptionHeader allows the setting a header option to configure
// the group.
func OptionHeader(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetHeader(i)
	}
}

// OptionHelp allows the setting a hint option to configure
// the group.
func OptionHelp(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetHelp(i)
	}
}

// OptionHint allows the setting a hint option to configure
// the group.
func OptionHint(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetHint(i)
	}
}

// OptionErr allows the setting a hint option to configure
// the group.
func OptionErr(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetErr(i)
	}
}

// OptionCommands allows the setting a commands option to configure
// the group.
func OptionCommands(i map[string]Command) HelpOption {
	return func(opt HelpOptions) {
		opt.SetCommands(i)
	}
}

// OptionFlags allows the setting a commands option to configure
// the group.
func OptionFlags(i []string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetFlags(i)
	}
}

// OptionFormat allows the setting a format option to configure
// the group.
func OptionFormat(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetFormat(i)
	}
}

// OptionColor allows the setting a color option to configure
// the group.
func OptionColor(i bool) HelpOption {
	return func(opt HelpOptions) {
		opt.SetColor(i)
	}
}

// OptionTemplate allows the setting a template option to configure
// the group.
func OptionTemplate(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetTemplate(i)
	}
}

// Func is the type of the function that is responsible for generating the
// help output when the CLI must show the general help text.
type Func func(...HelpOption) (string, error)

type nameHelp struct {
	Name     string
	Synopsis string
}

// BasicFunc generates some bashic help output that is usually good enough
// for most CLI applications.
func BasicFunc(name string) Func {
	return func(options ...HelpOption) (string, error) {
		opt := new(help)
		for _, option := range options {
			option(opt)
		}

		var buf bytes.Buffer
		writer := tabwriter.NewWriter(&buf, 24, 4, 4, ' ', 0)

		serialized := make([]nameHelp, len(opt.commands))

		var i int
		for k, v := range opt.commands {
			serialized[i] = nameHelp{
				Name:     k,
				Synopsis: v.Synopsis(),
			}
			i++
		}

		sort.Slice(serialized, func(i, j int) bool {
			return serialized[i].Name < serialized[j].Name
		})

		format := opt.format
		if strings.TrimSpace(format) == "" {
			format = HelpTemplateFormat
		}
		template := opt.template
		if strings.TrimSpace(template) == "" {
			template = BasicHelpTemplate
		}
		formatted := fmt.Sprintf(template, format)

		t := ui.NewTemplate(formatted,
			ui.OptionName("basic-help:"+name),
			ui.OptionColor(opt.color),
		)
		if err := t.Write(writer, struct {
			Name     string
			Header   string
			Hint     string
			Help     string
			Err      string
			Commands []nameHelp
			Flags    []string
		}{
			Name:     name,
			Header:   opt.header,
			Hint:     opt.hint,
			Help:     opt.help,
			Err:      opt.err,
			Commands: serialized,
			Flags:    opt.flags,
		}); err != nil {
			return "", errors.WithStack(err)
		}

		if err := writer.Flush(); err != nil {
			return "", errors.WithStack(err)
		}

		return strings.TrimSpace(buf.String()) + "\n", nil
	}
}
