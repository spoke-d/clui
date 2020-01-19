package args

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// FileSystem is an abstraction over the native filesystem
type FileSystem interface {

	// Stat returns a FileInfo describing the named file.
	Stat(string) (os.FileInfo, error)

	// Getwd returns a rooted path name corresponding to the current directory.
	Getwd() (string, error)
}

// Args describes command line arguments
type Args struct {
	fs FileSystem

	// All lists of all arguments in command line (not including the command itself)
	all []string
	// Completed lists of all completed arguments in command line,
	// If the last one is still being typed - no space after it,
	// it won't appear in this list of arguments.
	completed []string
	// Last argument in command line, the one being typed, if the last
	// character in the command line is a space, this argument will be empty,
	// otherwise this would be the last word.
	last string
	// LastCompleted is the last argument that was fully typed.
	// If the last character in the command line is space, this would be the
	// last word, otherwise, it would be the word before that.
	lastCompleted string
}

// ArgOptions represents a way to set optional values to a installer
// option.
// The ArgOptions shows what options are available to change.
type ArgOptions interface {
	SetFileSystem(FileSystem)
}

// ArgOption captures a tweak that can be applied to the Arg.
type ArgOption func(ArgOptions)

type args struct {
	fileSystem FileSystem
}

func (s *args) SetFileSystem(fs FileSystem) {
	s.fileSystem = fs
}

// OptionFileSystem allows the setting a filesystem option to configure
// the args.
func OptionFileSystem(fs FileSystem) ArgOption {
	return func(opt ArgOptions) {
		opt.SetFileSystem(fs)
	}
}

// New creates a Args from the command line arguments.
func New(line string, options ...ArgOption) *Args {
	opt := new(args)
	for _, option := range options {
		option(opt)
	}

	var (
		all       []string
		completed []string
	)
	parts := splitFields(line)
	if len(parts) > 0 {
		all = parts[1:]
		completed = removeLast(parts[1:])
	}
	return &Args{
		fs:            opt.fileSystem,
		all:           all,
		completed:     completed,
		last:          last(parts),
		lastCompleted: last(completed),
	}
}

// AllCommands returns all the commands.
func (a *Args) AllCommands() (c []string) {
	for _, v := range a.all {
		if !strings.HasPrefix(v, "-") {
			c = append(c, v)
		}
	}
	return
}

// CompletedCommands returns all the potential completed commands
func (a *Args) CompletedCommands() (c []string) {
	for _, v := range a.completed {
		if !strings.HasPrefix(v, "-") {
			c = append(c, v)
		}
	}
	return
}

// Directory gives the directory of the current written
// last argument if it represents a file name being written.
// in case that it is not, we fall back to the current directory.
func (a *Args) Directory() string {
	if info, err := a.fs.Stat(a.last); err == nil && info.IsDir() {
		return a.fixPathForm(a.last, a.last)
	}
	dir := filepath.Dir(a.last)
	if info, err := a.fs.Stat(dir); err != nil || !info.IsDir() {
		return "./"
	}
	return a.fixPathForm(a.last, dir)
}

// From captures a set of Args from a index.
func (a *Args) From(i int) *Args {
	if i > len(a.all) {
		i = len(a.all)
	}
	a.all = a.all[i:]

	if i > len(a.completed) {
		i = len(a.completed)
	}
	a.completed = a.completed[i:]
	return a
}

// Last returns the last argument
func (a *Args) Last() string {
	return a.last
}

// LastCompleted returns the last completed argument
func (a *Args) LastCompleted() string {
	return a.lastCompleted
}

// fixPathForm changes a file name to a relative name
func (a *Args) fixPathForm(last string, file string) string {
	// get working directory for relative name
	workDir, err := a.fs.Getwd()
	if err != nil {
		return file
	}

	abs, err := filepath.Abs(file)
	if err != nil {
		return file
	}

	// if last is absolute, return path as absolute
	if filepath.IsAbs(last) {
		return a.fixDirPath(abs)
	}

	rel, err := filepath.Rel(workDir, abs)
	if err != nil {
		return file
	}

	// fix ./ prefix of path
	if rel != "." && strings.HasPrefix(last, ".") {
		rel = "./" + rel
	}

	return a.fixDirPath(rel)
}

func (a *Args) fixDirPath(path string) string {
	info, err := a.fs.Stat(path)
	if err == nil && info.IsDir() && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

func splitFields(line string) []string {
	parts := strings.Fields(line)
	if len(line) > 0 && unicode.IsSpace(rune(line[len(line)-1])) {
		parts = append(parts, "")
	}
	parts = splitLastEqual(parts)
	return parts
}

func splitLastEqual(line []string) []string {
	if len(line) == 0 {
		return line
	}
	parts := strings.Split(line[len(line)-1], "=")
	return append(line[:len(line)-1], parts...)
}

func removeLast(a []string) []string {
	if len(a) > 0 {
		return a[:len(a)-1]
	}
	return a
}

func last(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[len(args)-1]
}
