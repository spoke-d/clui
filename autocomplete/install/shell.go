package install

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spoke-d/clui/autocomplete/fsys"
)

// Shell represents a shell for installing and uninstalling autocomplete
// commands in the correct shell files.
type Shell struct {
	fsys  fsys.FileSystem
	files []string
	cmdFn func(string, string) string
}

// Install attempts to install an autocomplete command into the correct file
// locations for that shell.
// Returns an error if the command already exists or if appending to a file
// is not possible.
func (b *Shell) Install(cmd, bin string) error {
	c := b.cmdFn(cmd, bin)
	for _, f := range b.files {
		if fileContains(b.fsys, f, c) {
			return fmt.Errorf("file already contains line: %q", c)
		}
		if err := appendToFile(b.fsys, f, c); err != nil {
			return err
		}
		// Only append to one file, otherwise it's just crazy.
		break
	}
	return nil
}

// Uninstall attempts to uninstall the autocomplete command from the shell
// files.
// Returns an error if it's unable to modify the profile file.
func (b *Shell) Uninstall(cmd, bin string) error {
	c := b.cmdFn(cmd, bin)
	for _, f := range b.files {
		if !fileContains(b.fsys, f, c) {
			continue
		}

		if err := removeFromFile(b.fsys, f, c); err != nil {
			return err
		}
	}
	return nil
}

// User is an abstraction around the current User.
type User interface {

	// HomeDir returns the home directory of the current User.
	HomeDir() string
}

// ShellOptions represents a way to set optional values to a installer
// option.
// The ShellOptions shows what options are available to change.
type ShellOptions interface {
	SetFileSystem(fsys.FileSystem)
	SetUser(User)
}

// ShellOption captures a tweak that can be applied to the Shell.
type ShellOption func(ShellOptions)

type shell struct {
	fileSystem fsys.FileSystem
	user       User
}

func (s *shell) SetFileSystem(fs fsys.FileSystem) {
	s.fileSystem = fs
}

func (s *shell) SetUser(u User) {
	s.user = u
}

// OptionFileSystem allows the setting a filesystem option to configure
// the shell.
func OptionFileSystem(fs fsys.FileSystem) ShellOption {
	return func(opt ShellOptions) {
		opt.SetFileSystem(fs)
	}
}

// OptionUser allows the setting a user option to configure the shell.
func OptionUser(u User) ShellOption {
	return func(opt ShellOptions) {
		opt.SetUser(u)
	}
}

// Bash creates a new Shell with the file locations and the correct bash
// complete command.
func Bash(options ...ShellOption) *Shell {
	opt := new(shell)
	for _, option := range options {
		option(opt)
	}

	var files []string

	for _, rc := range []string{
		".bashrc",
		".bash_profile",
		".bash_login",
		".profile",
	} {
		if path, ok := filePath(opt.fileSystem, opt.user, rc); ok {
			files = append(files, path)
		}
	}

	return &Shell{
		fsys:  opt.fileSystem,
		files: files,
		cmdFn: func(cmd, bin string) string {
			return fmt.Sprintf("complete -C %s %s", bin, cmd)
		},
	}
}

// Zsh creates a new Shell with the file locations and the correct zsh complete
// command.
func Zsh(options ...ShellOption) *Shell {
	opt := new(shell)
	for _, option := range options {
		option(opt)
	}

	var files []string

	for _, rc := range []string{
		".zshrc",
	} {
		if path, ok := filePath(opt.fileSystem, opt.user, rc); ok {
			files = append(files, path)
		}
	}

	return &Shell{
		fsys:  opt.fileSystem,
		files: files,
		cmdFn: func(cmd, bin string) string {
			return fmt.Sprintf("complete -o nospace -C %s %s", bin, cmd)
		},
	}
}

func filePath(fsys fsys.FileSystem, user User, file string) (string, bool) {
	path := filepath.Join(user.HomeDir(), file)
	return path, path != "" && fsys.Exists(path)
}

func fileContains(fsys fsys.FileSystem, name, pattern string) bool {
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

func appendToFile(fsys fsys.FileSystem, name, content string) error {
	f, err := fsys.OpenFile(name, os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(fmt.Sprintf("\n%s\n", content)))
	return err
}

func removeFromFile(fsys fsys.FileSystem, name, content string) error {
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

	return fsys.Remove(backupName)
}

func copyFile(fsys fsys.FileSystem, src, dst string) error {
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

func removeContentFromTmpFile(fsys fsys.FileSystem, name, content string) (string, error) {
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
