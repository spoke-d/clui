package group

import (
	"reflect"
	"sort"
	"strings"
	"testing"
	"testing/quick"

	"github.com/golang/mock/gomock"
	"github.com/spoke-d/clui/radix"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()

	t.Run("add", func(t *testing.T) {
		fn := guard(func(key string) bool {
			group := New()
			return group.Add(key, nil) == nil
		})
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get", func(t *testing.T) {
		fn := guard(func(key string) bool {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := NewMockCommand(ctrl)

			group := New()
			group.Add(key, cmd)

			c, ok := group.Get(key)
			return c == cmd && ok
		})
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("get other", func(t *testing.T) {
		fn := guard(func(key string) bool {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := NewMockCommand(ctrl)

			group := New()
			group.Add(key, cmd)

			_, ok := group.Get("#" + key)
			return !ok
		})
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("remove", func(t *testing.T) {
		fn := guard(func(key string) bool {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := NewMockCommand(ctrl)

			group := New()
			group.Add(key, cmd)

			c, err := group.Remove(key)
			return c == cmd && err == nil
		})
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("remove outer", func(t *testing.T) {
		fn := guard(func(key string) bool {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := NewMockCommand(ctrl)

			group := New()
			group.Add(key, cmd)

			_, err := group.Remove("#" + key)
			return strings.Contains(err.Error(), "no valid key found")
		})
		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})
}

func TestWalkPrefix(t *testing.T) {
	t.Parallel()

	t.Run("walk with nothing", func(t *testing.T) {
		var called bool

		group := New()
		group.WalkPrefix("a", func(s string, v radix.Value) bool {
			called = true
			return true
		})

		if expected, actual := false, called; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("walk", func(t *testing.T) {
		var found []string

		group := New()
		group.Add("abc", nil)
		group.Add("aaa", nil)
		group.Add("axx", nil)
		group.Add("xxc", nil)
		group.WalkPrefix("a", func(s string, v radix.Value) bool {
			found = append(found, s)
			return false
		})

		sort.Strings(found)

		if expected, actual := []string{"aaa", "abc", "axx"}, found; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("walk nested", func(t *testing.T) {
		var found []string

		group := New()
		group.Add("xxx", nil)
		group.Add("xxx aaa", nil)
		group.Add("xxx abc", nil)
		group.WalkPrefix("xxx a", func(s string, v radix.Value) bool {
			found = append(found, s)
			return false
		})

		sort.Strings(found)

		if expected, actual := []string{"xxx aaa", "xxx abc"}, found; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestLongestPrefix(t *testing.T) {
	t.Parallel()

	t.Run("longest with nothing", func(t *testing.T) {
		group := New()
		_, ok := group.LongestPrefix("a")

		if expected, actual := false, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("longest", func(t *testing.T) {
		group := New()
		group.Add("aaa", nil)
		group.Add("aaa aaa", nil)
		group.Add("aaaaaa", nil)
		found, ok := group.LongestPrefix("aaaa")

		if expected, actual := "aaa", found; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := true, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("longest nested", func(t *testing.T) {
		group := New()
		group.Add("aaa", nil)
		group.Add("aaa aaa", nil)
		group.Add("aaaaaa", nil)
		found, ok := group.LongestPrefix("aaa aa")

		if expected, actual := "aaa", found; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := true, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestNested(t *testing.T) {
	t.Parallel()

	t.Run("no nested commands", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := NewMockCommand(ctrl)

		group := New()
		group.Add("a", cmd)

		nested := group.Nested()
		if expected, actual := false, nested; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("nested commands", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := NewMockCommand(ctrl)

		group := New()
		group.Add("a a", cmd)

		nested := group.Nested()
		if expected, actual := true, nested; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestProcess(t *testing.T) {
	t.Parallel()

	t.Run("no nested commands", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := NewMockCommand(ctrl)

		group := New()
		group.Add("a", cmd)

		err := group.Process()
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})

	t.Run("nested", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := NewMockCommand(ctrl)

		var called bool
		var found string
		group := New(OptionPlaceHolder(func(s string) Command {
			called = true
			found = s
			return cmd
		}))
		group.Add("a b", nil)

		_, ok := group.Get("a")
		if expected, actual := false, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}

		err := group.Process()
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := true, called; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := "a", found; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}

		foundCmd, ok := group.Get("a")
		if expected, actual := cmd, foundCmd; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
		if expected, actual := true, ok; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestNormalizeKey(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Key    string
		Result string
	}{
		{
			Key:    "a",
			Result: "a",
		},
		{
			Key:    "a b",
			Result: "a b",
		},
		{
			Key:    "a b ",
			Result: "a b",
		},
		{
			Key:    " a   b    x v      ",
			Result: "a b x v",
		},
	}
	for _, v := range cases {
		t.Run(v.Key, func(t *testing.T) {
			res := normalizeKey(v.Key)
			if expected, actual := v.Result, res; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

func TestGetClosestName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Key      string
		Complete string
		Result   string
		Valid    bool
	}{
		{
			Key:      "test this",
			Complete: "",
			Result:   "",
			Valid:    false,
		},
		{
			Key:      "test this",
			Complete: "test this",
			Result:   "test this",
			Valid:    true,
		},
		{
			Key:      "test this",
			Complete: "test th",
			Result:   "test this",
			Valid:    true,
		},
		{
			Key:      "something",
			Complete: "foobar said nothing parent child",
			Result:   "",
			Valid:    false,
		},
	}

	for _, v := range cases {
		t.Run(v.Key, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := NewMockCommand(ctrl)

			group := New()
			err := group.Add(v.Key, cmd)
			if err != nil {
				t.Error(err)
			}

			res, ok := group.GetClosestName(v.Complete)
			if expected, actual := v.Result, res; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
			if expected, actual := v.Valid, ok; expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

func guard(fn func(string) bool) func(string) bool {
	return func(name string) bool {
		if name == "" {
			return true
		}
		return fn(strings.Replace(name, " ", "", -1))
	}
}
