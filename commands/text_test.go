package commands

import (
	"github.com/spoke-d/task/group"
	"testing"
	"testing/quick"
)

func TestTextCommand(t *testing.T) {
	t.Parallel()

	t.Run("text", func(t *testing.T) {
		f := func(a, b string) bool {
			text := NewText(a, b)

			return text.Help() == a && text.Synopsis() == b
		}
		if err := quick.Check(f, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("run", func(t *testing.T) {
		group := group.NewGroup()

		text := NewText("", "")
		text.Run(group)

		err := group.Run()

		if expected, actual := ErrShowHelp, err; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

}
