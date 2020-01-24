package fsys

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const mkdirAllMode = 0755

type LocalFileSystem struct{}

// NewLocalFileSystem yields a local disk filesystem.
func NewLocalFileSystem() LocalFileSystem {
	return LocalFileSystem{}
}

// Create takes a path, creates the file and then returns a File back that
// can be used. This returns an error if the file can not be created in
// some way.
func (LocalFileSystem) Create(path string) (File, error) {
	f, err := os.Create(path)
	return LocalFile{
		File:   f,
		Reader: f,
		Closer: f,
	}, err
}

// Open takes a path, opens a potential file and then returns a File if
// that file exists, otherwise it returns an error if the file wasn't found.
func (fs LocalFileSystem) Open(path string) (File, error) {
	f, err := os.Open(path)
	return fs.open(f, err)
}

// OpenFile takes a path, opens a potential file and then returns a File if
// that file exists, otherwise it returns an error if the file wasn't found.
func (fs LocalFileSystem) OpenFile(path string, flag int, perm os.FileMode) (File, error) {
	f, err := os.OpenFile(path, flag, perm)
	return fs.open(f, err)
}

// Exists takes a path and checks to see if the potential file exists or
// not.
// Note: If there is an error trying to read that file, it will return false
// even if the file already exists.
func (LocalFileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Remove takes a path, removes a potential file, if no file doesn't exist it
// will return not found.
func (LocalFileSystem) Remove(path string) error {
	return os.Remove(path)
}

func (fs LocalFileSystem) open(f *os.File, err error) (LocalFile, error) {
	if err != nil {
		if err == os.ErrNotExist {
			return LocalFile{}, errors.Wrap(err, "not found")
		}
		return LocalFile{}, errors.WithStack(err)
	}

	return LocalFile{
		File:   f,
		Reader: f,
		Closer: f,
	}, nil
}

// LocalFile is an abstraction for reading, writing and also closing a file.
type LocalFile struct {
	*os.File
	io.Reader
	io.Closer
}

func (f LocalFile) Read(p []byte) (int, error) {
	return f.Reader.Read(p)
}

func (f LocalFile) Close() error {
	return f.Closer.Close()
}

// Size returns the size of the file
func (f LocalFile) Size() int64 {
	fi, err := f.File.Stat()
	if err != nil {
		return -1
	}
	return fi.Size()
}

type deletingReleaser struct {
	path string
	r    Releaser
}

func (dr deletingReleaser) Release() error {
	// Remove before Release should be safe, and prevents a race.
	if err := os.Remove(dr.path); err != nil {
		return err
	}
	return dr.r.Release()
}

// multiCloser closes all underlying io.Closers.
// If an error is encountered, closings continue.
type multiCloser []io.Closer

func (c multiCloser) Close() error {
	var errs []error
	for _, closer := range c {
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return multiCloseError(errs)
	}
	return nil
}

type multiCloseError []error

func (e multiCloseError) Error() string {
	a := make([]string, len(e))
	for i, err := range e {
		a[i] = err.Error()
	}
	return strings.Join(a, "; ")
}
