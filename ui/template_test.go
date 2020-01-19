package ui

import "testing"

func TestTemplate(t *testing.T) {
	t.Parallel()

	t.Run("template", func(t *testing.T) {
		template := "{{.Name}}"

		temp := NewTemplate(template)
		result, err := temp.Render(struct {
			Name string
		}{
			Name: "Fred",
		})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := "Fred", result; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("template with failed template", func(t *testing.T) {
		template := "{{.Name"

		temp := NewTemplate(template, OptionName("template"))
		_, err := temp.Render(struct {
			Name string
		}{
			Name: "Fred",
		})
		if expected, actual := false, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})
}
