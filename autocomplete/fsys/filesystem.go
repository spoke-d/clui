package fsys

import (
	"io"
	"os"
)

// FileSystem is an abstraction over the native filesystem
type FileSystem interface {

	// Create takes a path, creates the file and then returns a File back that
	// can be used. This returns an error if the file can not be created in
	// some way.
	Create(string) (File, error)

	// Open takes a path, opens a potential file and then returns a File if
	// that file exists, otherwise it returns an error if the file wasn't found.
	Open(string) (File, error)

	// OpenFile takes a path, opens a potential file and then returns a File if
	// that file exists, otherwise it returns an error if the file wasn't found.
	OpenFile(path string, flag int, perm os.FileMode) (File, error)

	// Exists takes a path and checks to see if the potential file exists or
	// not.
	// Note: If there is an error trying to read that file, it will return false
	// even if the file already exists.
	Exists(string) bool

	// Remove takes a path and attempts to remove the file supplied.
	Remove(string) error
}

// File is an abstraction for reading, writing and also closing a file. These
// interfaces already exist, it's just a matter of composing them to be more
// usable by other components.
type File interface {
	io.Reader
	io.Writer
	io.Closer

	// Name returns the name of the file
	Name() string

	// Size returns the size of the file
	Size() int64

	// Sync attempts to sync the file with the underlying storage or errors if it
	// can't not succeed.
	Sync() error
}

// Releaser is returned by Lock calls.
type Releaser interface {

	// Release given lock or returns error upon failure.
	Release() error
}
