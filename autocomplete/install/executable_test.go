package install

import (
	"errors"
	"strings"
	"testing"
)

func TestExecutable(t *testing.T) {
	t.Parallel()

	t.Run("binary path", func(t *testing.T) {
		exec := OSExecutable{}
		path, err := exec.getBinaryPath(func() (string, error) {
			return "foo/bar/baz", nil
		})
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
		if expected, actual := "/foo/bar/baz", path; !strings.Contains(actual, expected) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("binary path error", func(t *testing.T) {
		exec := OSExecutable{}
		_, err := exec.getBinaryPath(func() (string, error) {
			return "", errors.New("fail")
		})
		if expected, actual := "fail", err.Error(); expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})
}
