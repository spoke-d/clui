package ui

import (
	"bytes"
	"fmt"
	"testing"
	"testing/quick"
)

func TestBasicUI(t *testing.T) {
	t.Parallel()

	t.Run("info", func(t *testing.T) {
		f := func(s string) bool {
			var buf bytes.Buffer

			ui := NewBasicUI(nil, &buf, nil)
			ui.Info(s)

			return buf.String() == fmt.Sprintf("%s\n", s)
		}
		if err := quick.Check(f, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := func(s string) bool {
			var buf bytes.Buffer

			ui := NewBasicUI(nil, nil, &buf)
			ui.Error(s)

			return buf.String() == fmt.Sprintf("%s\n", s)
		}
		if err := quick.Check(f, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("output", func(t *testing.T) {
		var buf bytes.Buffer

		ui := NewBasicUI(nil, &buf, nil)
		ui.Output(NewTemplate("{{.Name}}"), struct {
			Name string
		}{
			Name: "Fred",
		})
		if expected, actual := "Fred\n", buf.String(); expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
