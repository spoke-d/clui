package install_test

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spoke-d/clui/install"
	"github.com/spoke-d/clui/install/mocks"
)

func TestInstaller(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	regex := regexp.MustCompile("complete.*-C .*/install.test xxx")

	fs := mocks.NewMockFileSystem(ctrl)
	file := mocks.NewMockFile(ctrl)

	fsExp := fs.EXPECT()
	fExp := file.EXPECT()

	fsExp.Exists(filepath.Join(u.HomeDir, ".bashrc")).Return(false)
	fsExp.Exists(filepath.Join(u.HomeDir, ".bash_profile")).Return(false)
	fsExp.Exists(filepath.Join(u.HomeDir, ".bash_login")).Return(false)
	fsExp.Exists(filepath.Join(u.HomeDir, ".profile")).Return(true)
	fsExp.Exists(filepath.Join(u.HomeDir, ".zshrc")).Return(true)

	for _, v := range []string{".profile", ".zshrc"} {
		fsExp.Open(filepath.Join(u.HomeDir, v)).Return(file, nil)
		fExp.Read(gomock.Any()).Return(0, io.EOF)
		fExp.Close().Return(nil)
		fsExp.OpenFile(filepath.Join(u.HomeDir, v), os.O_RDWR|os.O_APPEND, os.FileMode(0755)).Return(file, nil)
		fExp.Write(regexpMatcher{regex}).Return(0, nil)
		fExp.Close().Return(nil)
	}

	in, err := install.New(fs)
	if err != nil {
		t.Fatal(err)
	}

	if err := in.Install("xxx"); err != nil {
		t.Error(err)
	}
}

type regexpMatcher struct {
	regex *regexp.Regexp
}

func (m regexpMatcher) Matches(x interface{}) bool {
	if b, ok := x.([]byte); ok {
		return m.regex.MatchString(string(b))
	}
	return false
}

func (m regexpMatcher) String() string {
	return "want regexp " + m.regex.String()
}
