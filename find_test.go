package clui

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spoke-d/clui/group"
)

func TestFindChildren(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		group := group.New()
		children, err := FindChildren(group, "")
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := map[string]Command{}, children; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("empty group", func(t *testing.T) {
		group := group.New()
		children, err := FindChildren(group, "foo bar")
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := map[string]Command{}, children; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("group without command", func(t *testing.T) {
		group := group.New()
		group.Add("foo bar", nil)

		_, err := FindChildren(group, "foo")
		if expected, actual := "not found: \"foo bar\"", err.Error(); expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})

	t.Run("group", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := NewMockCommand(ctrl)

		group := group.New()
		group.Add("foo bar", cmd)

		children, err := FindChildren(group, "foo")
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := map[string]Command{"foo bar": cmd}, children; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
