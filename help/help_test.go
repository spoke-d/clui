package help

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestBasicFunc(t *testing.T) {
	t.Parallel()

	t.Run("no commands", func(t *testing.T) {
		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionShowHelp(true),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := `
Usage: foo [--version] [--help] [--debug] <command> [<args>]

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`[1:]
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("commands", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmdVersion := NewMockCommand(ctrl)
		cmdVersion.EXPECT().Synopsis().Return("returns the version")

		cmdFoo := NewMockCommand(ctrl)
		cmdFoo.EXPECT().Synopsis().Return("foo command")

		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionCommands(map[string]Command{
				"version":         cmdVersion,
				"foo bar baz xxx": cmdFoo,
			}),
			OptionShowHelp(true),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
Usage: foo [--version] [--help] [--debug] <command> [<args>]

Available commands:

    foo bar baz xxx     foo command
    version             returns the version

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`) + "\n"
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("commands with format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmdVersion := NewMockCommand(ctrl)
		cmdVersion.EXPECT().Synopsis().Return("returns the version")

		cmdFoo := NewMockCommand(ctrl)
		cmdFoo.EXPECT().Synopsis().Return("foo command")

		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionCommands(map[string]Command{
				"version":         cmdVersion,
				"foo bar baz xxx": cmdFoo,
			}),
			OptionFormat("{{.Name}}"),
			OptionShowHelp(true),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
Usage: foo [--version] [--help] [--debug] <command> [<args>]

Available commands:

foo bar baz xxx
version

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`) + "\n"
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmdVersion := NewMockCommand(ctrl)
		cmdVersion.EXPECT().Synopsis().Return("returns the version")

		cmdFoo := NewMockCommand(ctrl)
		cmdFoo.EXPECT().Synopsis().Return("foo command")

		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionCommands(map[string]Command{
				"version":         cmdVersion,
				"foo bar baz xxx": cmdFoo,
			}),
			OptionHeader("**HEADER**"),
			OptionShowHelp(true),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
**HEADER**

Usage: foo [--version] [--help] [--debug] <command> [<args>]

Available commands:

    foo bar baz xxx     foo command
    version             returns the version

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`) + "\n"
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("hint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmdVersion := NewMockCommand(ctrl)
		cmdVersion.EXPECT().Synopsis().Return("returns the version")

		cmdFoo := NewMockCommand(ctrl)
		cmdFoo.EXPECT().Synopsis().Return("foo command")

		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionCommands(map[string]Command{
				"version":         cmdVersion,
				"foo bar baz xxx": cmdFoo,
			}),
			OptionHeader("**HEADER**"),
			OptionHint("foo"),
			OptionShowHelp(true),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
**HEADER**

Did you mean?
        foo

Usage: foo [--version] [--help] [--debug] <command> [<args>]

Available commands:

    foo bar baz xxx     foo command
    version             returns the version

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`) + "\n"
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("hint without help", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmdVersion := NewMockCommand(ctrl)
		cmdVersion.EXPECT().Synopsis().Return("returns the version")

		cmdFoo := NewMockCommand(ctrl)
		cmdFoo.EXPECT().Synopsis().Return("foo command")

		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionCommands(map[string]Command{
				"version":         cmdVersion,
				"foo bar baz xxx": cmdFoo,
			}),
			OptionHeader("**HEADER**"),
			OptionHint("foo"),
			OptionShowHelp(false),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
**HEADER**

Did you mean?
        foo

`) + "\n"
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestComandFunc(t *testing.T) {
	t.Parallel()

	t.Run("no commands", func(t *testing.T) {
		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionErr("something went wrong"),
			OptionTemplate(CommandHelpTemplate),
			OptionShowHelp(true),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := `
Found some issues:

    something went wrong

See foo --help for more information.

Usage:

    foo

Description:
    

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`[1:]
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("no commands without help", func(t *testing.T) {
		helpFn := BasicFunc("foo")
		result, err := helpFn(
			OptionErr("something went wrong"),
			OptionTemplate(CommandHelpTemplate),
			OptionShowHelp(false),
		)

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := `
Found some issues:

    something went wrong

See foo --help for more information.
`[1:]
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
