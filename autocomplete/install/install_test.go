package install

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	fsys "github.com/spoke-d/clui/autocomplete/fsys"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("binary path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exec := NewMockExecutable(ctrl)
		exec.EXPECT().BinaryPath().Return("", errors.New("fail"))

		_, err := New(OptionExecutable(exec))
		if expected, actual := "fail", err.Error(); expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})
}

func TestInstall(t *testing.T) {
	t.Parallel()

	t.Run("install", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exec := NewMockExecutable(ctrl)
		exec.EXPECT().BinaryPath().Return("/home", nil)

		existingLine := "complete -C /home yyy"
		newLine := "complete -C /home xxx"

		file := NewMockFile(ctrl)
		fs := NewMockFileSystem(ctrl)

		// FileContains
		gomock.InOrder(
			fs.EXPECT().Open(".profile").Return(file, nil),
			file.EXPECT().Read(gomock.Any()).Return(len(existingLine), io.EOF).Do(func(b []byte) {
				for i := 0; i < len(existingLine); i++ {
					b[i] = existingLine[i]
				}
			}),
			file.EXPECT().Close(),
		)

		// AppendToFile
		gomock.InOrder(
			fs.EXPECT().OpenFile(".profile", gomock.Any(), gomock.Any()).Return(file, nil),
			file.EXPECT().Write([]byte(fmt.Sprintf("\n%s\n", newLine))).Return(0, nil),
			file.EXPECT().Close(),
		)

		in, err := New(OptionExecutable(exec), OptionShell(StubShell(fs)))
		if err != nil {
			t.Fail()
		}

		err = in.Install("xxx")
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})

	t.Run("install fails already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exec := NewMockExecutable(ctrl)
		exec.EXPECT().BinaryPath().Return("/home", nil)

		existingLine := "complete -C /home xxx"

		file := NewMockFile(ctrl)
		fs := NewMockFileSystem(ctrl)

		// FileContains
		gomock.InOrder(
			fs.EXPECT().Open(".profile").Return(file, nil),
			file.EXPECT().Read(gomock.Any()).Return(len(existingLine), io.EOF).Do(func(b []byte) {
				for i := 0; i < len(existingLine); i++ {
					b[i] = existingLine[i]
				}
			}),
			file.EXPECT().Close(),
		)

		in, err := New(OptionExecutable(exec), OptionShell(StubShell(fs)))
		if err != nil {
			t.Fail()
		}

		err = in.Install("xxx")
		if expected, actual := `file already contains line: "complete -C /home xxx"`, err.Error(); expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})

	t.Run("install fails appending file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exec := NewMockExecutable(ctrl)
		exec.EXPECT().BinaryPath().Return("/home", nil)

		existingLine := "complete -C /home yyy"

		file := NewMockFile(ctrl)
		fs := NewMockFileSystem(ctrl)

		// FileContains
		gomock.InOrder(
			fs.EXPECT().Open(".profile").Return(file, nil),
			file.EXPECT().Read(gomock.Any()).Return(len(existingLine), io.EOF).Do(func(b []byte) {
				for i := 0; i < len(existingLine); i++ {
					b[i] = existingLine[i]
				}
			}),
			file.EXPECT().Close(),
		)

		// AppendToFile
		gomock.InOrder(
			fs.EXPECT().OpenFile(".profile", gomock.Any(), gomock.Any()).Return(file, errors.New("fail")),
		)

		in, err := New(OptionExecutable(exec), OptionShell(StubShell(fs)))
		if err != nil {
			t.Fail()
		}

		err = in.Install("xxx")
		if expected, actual := "fail", err.Error(); expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})
}

func TestUninstall(t *testing.T) {
	t.Parallel()

	t.Run("uninstall", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		exec := NewMockExecutable(ctrl)
		exec.EXPECT().BinaryPath().Return("/home", nil)

		existingLine := "complete -C /home xxx"

		file := NewMockFile(ctrl)
		backupFile := NewMockFile(ctrl)
		tmpFile := NewMockFile(ctrl)
		fs := NewMockFileSystem(ctrl)

		// FileContains
		gomock.InOrder(
			fs.EXPECT().Open(".profile").Return(file, nil),
			file.EXPECT().Read(gomock.Any()).Return(len(existingLine), io.EOF).Do(func(b []byte) {
				for i := 0; i < len(existingLine); i++ {
					b[i] = existingLine[i]
				}
			}),
			file.EXPECT().Close(),
		)

		// CopyFile
		gomock.InOrder(
			fs.EXPECT().Open(".profile").Return(file, nil),
			fs.EXPECT().Create(".profile.bck").Return(backupFile, nil),
			file.EXPECT().Read(gomock.Any()).Return(len(existingLine), io.EOF).Do(func(b []byte) {
				for i := 0; i < len(existingLine); i++ {
					b[i] = existingLine[i]
				}
			}),
			backupFile.EXPECT().Write([]byte(fmt.Sprintf("%s", existingLine))).Return(len(existingLine), nil),
			backupFile.EXPECT().Close(),
			file.EXPECT().Close(),
		)

		// RemoveContentFromTmpFile
		gomock.InOrder(
			fs.EXPECT().Open(".profile").Return(file, nil),
			file.EXPECT().Read(gomock.Any()).Return(len(existingLine), io.EOF).Do(func(b []byte) {
				for i := 0; i < len(existingLine); i++ {
					b[i] = existingLine[i]
				}
			}),
			file.EXPECT().Close().Return(nil),
		)

		// CopyFile
		gomock.InOrder(
			fs.EXPECT().Open(gomock.Any()).Return(tmpFile, nil),
			fs.EXPECT().Create(".profile").Return(file, nil),
			tmpFile.EXPECT().Read(gomock.Any()).Return(len(existingLine), io.EOF).Do(func(b []byte) {
				for i := 0; i < len(existingLine); i++ {
					b[i] = existingLine[i]
				}
			}),
			file.EXPECT().Write([]byte(fmt.Sprintf("%s", existingLine))).Return(len(existingLine), nil),
			file.EXPECT().Close(),
			tmpFile.EXPECT().Close(),
		)

		// Remove
		gomock.InOrder(
			fs.EXPECT().Remove(".profile.bck").Return(nil),
		)

		in, err := New(OptionExecutable(exec), OptionShell(StubShell(fs)))
		if err != nil {
			t.Fail()
		}

		err = in.Uninstall("xxx")
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})
}

func StubShell(fsys fsys.FileSystem) *Shell {
	return &Shell{
		fsys:  fsys,
		files: []string{".profile"},
		cmdFn: func(cmd, bin string) string {
			return fmt.Sprintf("complete -C %s %s", bin, cmd)
		},
	}
}
