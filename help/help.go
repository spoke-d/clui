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
	SetCommands(map[string]Command)
	SetFormat(string)
	SetColor(bool)
}

// HelpOption captures a tweak that can be applied to the Help.
type HelpOption func(HelpOptions)

// HelpOptions defines options for overriding help rendering.
type help struct {
	header   string
	hint     string
	commands map[string]Command
	format   string
	color    bool
}

func (s *help) SetHeader(p string) {
	s.header = p
}

func (s *help) SetHint(p string) {
	s.hint = p
}

func (s *help) SetCommands(p map[string]Command) {
	s.commands = p
}

func (s *help) SetFormat(p string) {
	s.format = p
}

func (s *help) SetColor(p bool) {
	s.color = p
}

// OptionHeader allows the setting a header option to configure
// the group.
func OptionHeader(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetHeader(i)
	}
}

// OptionHint allows the setting a hint option to configure
// the group.
func OptionHint(i string) HelpOption {
	return func(opt HelpOptions) {
		opt.SetHint(i)
	}
}

// OptionCommands allows the setting a commands option to configure
// the group.
func OptionCommands(i map[string]Command) HelpOption {
	return func(opt HelpOptions) {
		opt.SetCommands(i)
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
		formatted := fmt.Sprintf(HelpTemplate, format)

		t := ui.NewTemplate(formatted,
			ui.OptionName("basic-help"),
			ui.OptionColor(opt.color),
		)
		if err := t.Write(writer, struct {
			Name     string
			Header   string
			Hint     string
			Commands []nameHelp
		}{
			Name:     name,
			Header:   opt.header,
			Hint:     opt.hint,
			Commands: serialized,
		}); err != nil {
			return "", errors.WithStack(err)
		}

		if err := writer.Flush(); err != nil {
			return "", errors.WithStack(err)
		}

		return strings.TrimSpace(buf.String()) + "\n", nil
	}
}
