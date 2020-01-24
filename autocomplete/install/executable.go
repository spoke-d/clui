package install

import (
	"os"
	"path/filepath"
)

// ExecutableFn is used to return the current binary path, so that we know
// which application to call when inserting the completion command.
type ExecutableFn func() (string, error)

// OSExecutable is a wrapper around the OS package to provide the binary path
// for the current executable.
type OSExecutable struct{}

// BinaryPath returns the binary path as an absolute file path.
func (e OSExecutable) BinaryPath() (string, error) {
	return e.getBinaryPath(os.Executable)
}

func (e OSExecutable) getBinaryPath(fn ExecutableFn) (string, error) {
	path, err := fn()
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}
