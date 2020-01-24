package clui

import (
	"reflect"
	"testing"

	"github.com/spoke-d/clui/group"
)

func TestGlobalArgs(t *testing.T) {
	t.Parallel()

	t.Run("no args", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})

	t.Run("help args (short)", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"-h"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.Help(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("help args (long)", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"--help"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.Help(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("version args (short)", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"-v"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.Version(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("version args (long)", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"--version"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.Version(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("debug args", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"--debug"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.Debug(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("no-color args", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"--no-color"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.RequiresNoColor(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("autocomplete install", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"--autocomplete-install"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.RequiresInstall(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("autocomplete uninstall", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"--autocomplete-uninstall"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, args.RequiresUninstall(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("autocomplete install and uninstall", func(t *testing.T) {
		group := group.New()

		args := NewGlobalArgs(group)
		err := args.Process([]string{"--autocomplete-uninstall", "--autocomplete-install"})
		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})

	//

	t.Run("command args", func(t *testing.T) {
		group := group.New()
		group.Add("a b", nil)

		args := NewGlobalArgs(group)
		err := args.Process([]string{"a", "b", "c"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := "a b", args.SubCommand(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string{"c"}, args.SubCommandArgs(); !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string(nil), args.CommandFlags(); !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("command args with flags", func(t *testing.T) {
		group := group.New()
		group.Add("a b", nil)

		args := NewGlobalArgs(group)
		err := args.Process([]string{"a", "b", "c", "--flag"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := "a b", args.SubCommand(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string{"c", "--flag"}, args.SubCommandArgs(); !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string(nil), args.CommandFlags(); !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("command args with flags", func(t *testing.T) {
		group := group.New()
		group.Add("a b", nil)

		args := NewGlobalArgs(group)
		err := args.Process([]string{"a", "b", "--flag", "c"})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := "a b", args.SubCommand(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string{"--flag", "c"}, args.SubCommandArgs(); !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string(nil), args.CommandFlags(); !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
