package args

import (
	"fmt"
	"os"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line          string
		completed     string
		last          string
		lastCompleted string
	}{
		{
			line:          "a b c",
			completed:     "b",
			last:          "c",
			lastCompleted: "b",
		},
		{
			line:          "a b ",
			completed:     "b",
			last:          "",
			lastCompleted: "b",
		},
		{
			line:          "",
			completed:     "",
			last:          "",
			lastCompleted: "",
		},
		{
			line:          "a",
			completed:     "",
			last:          "a",
			lastCompleted: "",
		},
		{
			line:          "a ",
			completed:     "",
			last:          "",
			lastCompleted: "",
		},
		{
			line:          "a -echo",
			completed:     "",
			last:          "-echo",
			lastCompleted: "",
		},
		{
			line:          "a -echo ",
			completed:     "-echo",
			last:          "",
			lastCompleted: "-echo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {

			a := New(tt.line)

			if got, want := strings.Join(a.completed, " "), tt.completed; got != want {
				t.Errorf("%s failed: Completed = %q, want %q", t.Name(), got, want)
			}
			if got, want := a.last, tt.last; got != want {
				t.Errorf("Last = %q, want %q", got, want)
			}
			if got, want := a.lastCompleted, tt.lastCompleted; got != want {
				t.Errorf("%s failed: LastCompleted = %q, want %q", t.Name(), got, want)
			}
		})
	}
}

func TestArgs_AllCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line     string
		commands string
	}{
		{
			line:     "a b c",
			commands: "b c",
		},
		{
			line:     "a b",
			commands: "b",
		},
		{
			line:     "a --foo b -bar",
			commands: "b",
		},
		{
			line:     "a --foo b -bar c --baz",
			commands: "b c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			a := New(tt.line)

			if got, want := strings.Join(a.AllCommands(), " "), tt.commands; got != want {
				t.Errorf("%s failed: Completed = %q, want %q", t.Name(), got, want)
			}
		})
	}
}

func TestArgs_CompletedCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line     string
		commands string
	}{
		{
			line:     "a b c",
			commands: "b",
		},
		{
			line:     "a b",
			commands: "",
		},
		{
			line:     "a --foo b -bar",
			commands: "b",
		},
		{
			line:     "a --foo b -bar c --baz",
			commands: "b c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {

			a := New(tt.line)

			if got, want := strings.Join(a.CompletedCommands(), " "), tt.commands; got != want {
				t.Errorf("%s failed: Completed = %q, want %q", t.Name(), got, want)
			}
		})
	}
}

func TestArgs_From(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line         string
		from         int
		newLine      string
		newCompleted string
	}{
		{
			line:         "a b c",
			from:         0,
			newLine:      "b c",
			newCompleted: "b",
		},
		{
			line:         "a b c",
			from:         1,
			newLine:      "c",
			newCompleted: "",
		},
		{
			line:         "a b c",
			from:         2,
			newLine:      "",
			newCompleted: "",
		},
		{
			line:         "a b c",
			from:         3,
			newLine:      "",
			newCompleted: "",
		},
		{
			line:         "a b c ",
			from:         0,
			newLine:      "b c ",
			newCompleted: "b c",
		},
		{
			line:         "a b c ",
			from:         1,
			newLine:      "c ",
			newCompleted: "c",
		},
		{
			line:         "a b c ",
			from:         2,
			newLine:      "",
			newCompleted: "",
		},
		{
			line:         "",
			from:         0,
			newLine:      "",
			newCompleted: "",
		},
		{
			line:         "",
			from:         1,
			newLine:      "",
			newCompleted: "",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%d", tt.line, tt.from), func(t *testing.T) {

			a := New(tt.line)
			n := a.From(tt.from)

			if got, want := strings.Join(n.all, " "), tt.newLine; got != want {
				t.Errorf("%s failed: all = %q, want %q", t.Name(), got, want)
			}
			if got, want := strings.Join(n.completed, " "), tt.newCompleted; got != want {
				t.Errorf("%s failed: completed = %q, want %q", t.Name(), got, want)
			}
		})
	}
}

func TestArgs_Directory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line      string
		directory string
	}{
		{
			line:      "a b c",
			directory: "./",
		},
		{
			line:      "a b c /tmp ",
			directory: "./",
		},
		{
			line:      "a b c ./",
			directory: "./",
		},
		{
			line:      "a b c ./di",
			directory: "./",
		},
		{
			line:      "a b c ./dir ",
			directory: "./",
		},
		{
			line:      "a b c ./di",
			directory: "./",
		},
		{
			line:      "a b c ./a.txt",
			directory: "./",
		},
		{
			line:      "a b c ./a.txt/x",
			directory: "./",
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			info := NewMockFileInfo(ctrl)
			info.EXPECT().IsDir().Return(false).AnyTimes()

			fs := NewMockFileSystem(ctrl)
			fs.EXPECT().Stat(gomock.Any()).Return(info, nil).AnyTimes()

			a := New(tt.line, OptionFileSystem(fs))

			if got, want := a.Directory(), tt.directory; got != want {
				t.Errorf("%s failed: directory = %q, want %q", t.Name(), got, want)
			}
		})
	}
}

func TestArgs_DirectoryFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line      string
		directory string
	}{
		{
			line:      "a b c /tmp",
			directory: "/tmp/",
		},
		{
			line:      "a b c ./tmp",
			directory: "./tmp/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wd, err := os.Getwd()
			if err != nil {
				t.Error(err)
			}

			info := NewMockFileInfo(ctrl)
			info.EXPECT().IsDir().Return(true).AnyTimes()

			fs := NewMockFileSystem(ctrl)
			fs.EXPECT().Stat(gomock.Any()).Return(info, nil).AnyTimes()
			fs.EXPECT().Getwd().Return(wd, nil)

			a := New(tt.line, OptionFileSystem(fs))

			if got, want := a.Directory(), tt.directory; got != want {
				t.Errorf("%s failed: directory = %q, want %q", t.Name(), got, want)
			}
		})
	}
}
