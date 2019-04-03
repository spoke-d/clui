package install

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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

type installer interface {
	Install(string, string) error
	Uninstall(string, string) error
}

// ExecutableFn is used to return the current binary path, so that we know
// which application to call when inserting the completion command.
type ExecutableFn func() (string, error)

// Installer represents a series of different shell installers. Calling
// Install and Uninstall to add and remove the correct shell installer commands
// retrospectively.
type Installer struct {
	installers []installer
	binaryPath string
}

// New creates a new installer with the correct dependencies and sane
// defaults.
// Returns an error if it can't locate the correct binary path.
func New(fsys FileSystem) (*Installer, error) {
	path, err := getBinaryPath(os.Executable)
	if err != nil {
		return nil, err
	}

	in := &Installer{
		installers: []installer{
			newBash(fsys),
			newZsh(fsys),
		},
		binaryPath: path,
	}
	return in, nil
}

// Install takes a command argument and installs the command for the shell
func (i *Installer) Install(cmd string) error {
	for _, in := range i.installers {
		if err := in.Install(cmd, i.binaryPath); err != nil {
			return err
		}
	}
	return nil
}

// Uninstall removes the auto-complete installer shell
func (i *Installer) Uninstall(cmd string) error {
	for _, in := range i.installers {
		if err := in.Uninstall(cmd, i.binaryPath); err != nil {
			return err
		}
	}
	return nil
}

// NopInstaller can be used as a default nop installer.
type NopInstaller struct{}

// NewNop creates an installer that performs no operations
func NewNop() *NopInstaller {
	return &NopInstaller{}
}

// Install takes a command argument and installs the command for the shell
func (*NopInstaller) Install(cmd string) error { return nil }

// Uninstall removes the auto-complete installer shell
func (*NopInstaller) Uninstall(cmd string) error { return nil }

func getBinaryPath(fn ExecutableFn) (string, error) {
	path, err := fn()
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

func filePath(fsys FileSystem, file string) (string, bool) {
	u, err := user.Current()
	if err != nil {
		return "", false
	}

	path := filepath.Join(u.HomeDir, file)
	return path, path != "" && fsys.Exists(path)
}

func fileContains(fsys FileSystem, name, pattern string) bool {
	f, err := fsys.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), pattern) {
			return true
		}
	}
	return false
}

func appendToFile(fsys FileSystem, name, content string) error {
	f, err := fsys.OpenFile(name, os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	fmt.Println(content)
	_, err = f.Write([]byte(fmt.Sprintf("\n%s\n", content)))
	return err
}

func removeFromFile(fsys FileSystem, name, content string) error {
	backupName := fmt.Sprintf("%s.bck", name)
	if err := copyFile(fsys, name, backupName); err != nil {
		return err
	}

	tmp, err := removeContentFromTmpFile(fsys, name, content)
	if err != nil {
		return err
	}

	if err := copyFile(fsys, tmp, name); err != nil {
		return err
	}

	return os.Remove(backupName)
}

func copyFile(fsys FileSystem, src, dst string) error {
	in, err := fsys.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := fsys.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func removeContentFromTmpFile(fsys FileSystem, name, content string) (string, error) {
	file, err := fsys.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	tmpFile, err := ioutil.TempFile("/tmp", "complete-")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, content) {
			continue
		}
		if _, err := tmpFile.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return "", err
		}
	}

	return tmpFile.Name(), nil
}
