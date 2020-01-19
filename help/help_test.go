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
		result, err := helpFn()

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := `
Usage: foo [--version] [--help] <command> [<args>]
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
		result, err := helpFn(OptionCommands(map[string]Command{
			"version":         cmdVersion,
			"foo bar baz xxx": cmdFoo,
		}))

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
Usage: foo [--version] [--help] <command> [<args>]

Available commands are:

foo bar baz xxx         foo command
version                 returns the version
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
		result, err := helpFn(OptionCommands(map[string]Command{
			"version":         cmdVersion,
			"foo bar baz xxx": cmdFoo,
		}), OptionFormat("{{.Name}}"))

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
Usage: foo [--version] [--help] <command> [<args>]

Available commands are:

foo bar baz xxx
version
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
		result, err := helpFn(OptionCommands(map[string]Command{
			"version":         cmdVersion,
			"foo bar baz xxx": cmdFoo,
		}), OptionHeader("**HEADER**"))

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
**HEADER**

Usage: foo [--version] [--help] <command> [<args>]

Available commands are:

foo bar baz xxx         foo command
version                 returns the version
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
		result, err := helpFn(OptionCommands(map[string]Command{
			"version":         cmdVersion,
			"foo bar baz xxx": cmdFoo,
		}), OptionHeader("**HEADER**"), OptionHint("foo"))

		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		required := strings.TrimSpace(`
**HEADER**

Did you mean?
        foo

Usage: foo [--version] [--help] <command> [<args>]

Available commands are:

foo bar baz xxx         foo command
version                 returns the version
`) + "\n"
		if expected, actual := required, result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
