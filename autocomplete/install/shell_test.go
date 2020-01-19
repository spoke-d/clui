package install

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestBash(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fs := NewMockFileSystem(ctrl)
		u := NewMockUser(ctrl)

		u.EXPECT().HomeDir().Return("/home/test").Times(4)

		gomock.InOrder(
			fs.EXPECT().Exists("/home/test/.bashrc").Return(false),
			fs.EXPECT().Exists("/home/test/.bash_profile").Return(true),
			fs.EXPECT().Exists("/home/test/.bash_login").Return(false),
			fs.EXPECT().Exists("/home/test/.profile").Return(true),
		)

		shell := Bash(OptionFileSystem(fs), OptionUser(u))
		cmd := shell.cmdFn("xxx", "file.bin")
		if expected, actual := "complete -C file.bin xxx", cmd; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestZsh(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fs := NewMockFileSystem(ctrl)
		u := NewMockUser(ctrl)

		u.EXPECT().HomeDir().Return("/home/test")

		gomock.InOrder(
			fs.EXPECT().Exists("/home/test/.zshrc").Return(true),
		)

		shell := Zsh(OptionFileSystem(fs), OptionUser(u))
		cmd := shell.cmdFn("xxx", "file.bin")
		if expected, actual := "complete -o nospace -C file.bin xxx", cmd; expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
