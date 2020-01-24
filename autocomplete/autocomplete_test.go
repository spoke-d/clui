package autocomplete

import (
	"flag"
	reflect "reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/radix"
)

func TestAutoCompleteCommand(t *testing.T) {
	t.Parallel()

	t.Run("complete empty line", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		group := NewMockGroup(ctrl)

		ac := New(OptionGroup(group))
		_, ok := ac.Complete("")
		if expected, actual := false, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("complete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		group := NewMockGroup(ctrl)
		group.EXPECT().WalkPrefix("test foo", gomock.Any()).Do(func(s string, fn func(s string, cmd radix.Value) bool) {
			fn(s, NewMockCommand(ctrl))
		})

		ac := New(OptionGroup(group))
		matches, ok := ac.Complete("clui test foo")
		if expected, actual := true, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string{"foo"}, matches; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("complete with trailing space", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		group := NewMockGroup(ctrl)
		group.EXPECT().WalkPrefix("test foo ", gomock.Any()).Do(func(s string, fn func(s string, cmd radix.Value) bool) {
			fn("bar", NewMockCommand(ctrl))
		})

		ac := New(OptionGroup(group))
		matches, ok := ac.Complete("clui test foo ")
		if expected, actual := true, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string{"bar"}, matches; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestAutoCompleteFlagset(t *testing.T) {
	t.Parallel()

	t.Run("complete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		flagSet := flagset.New("test", flag.ContinueOnError)
		flagSet.String("bar", "false", "some usage pattern here")

		cmd := NewMockCommand(ctrl)
		cmd.EXPECT().FlagSet().Return(flagSet)

		group := NewMockGroup(ctrl)
		group.EXPECT().WalkPrefix("test foo", gomock.Any()).Do(func(s string, fn func(s string, cmd radix.Value) bool) {
			fn(s, cmd)
		})

		ac := New(OptionGroup(group))
		matches, ok := ac.Complete("clui test foo --bar")
		if expected, actual := true, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string{"--bar"}, matches; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("complete with multiple", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		flagSet := flagset.New("test", flag.ContinueOnError)
		flagSet.String("bar", "false", "some usage pattern here")
		flagSet.String("baz", "false", "some usage pattern here")

		cmd := NewMockCommand(ctrl)
		cmd.EXPECT().FlagSet().Return(flagSet)

		group := NewMockGroup(ctrl)
		group.EXPECT().WalkPrefix("test foo", gomock.Any()).Do(func(s string, fn func(s string, cmd radix.Value) bool) {
			fn(s, cmd)
		})

		ac := New(OptionGroup(group))
		matches, ok := ac.Complete("clui test foo --ba")
		if expected, actual := true, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := []string{"--bar", "--baz"}, matches; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
