package install

// Executable defines an abstraction over the current binary executable.
type Executable interface {
	// BinaryPath returns the current binary path of the executable.
	BinaryPath() (string, error)
}

// Installer represents a series of different shell installers. Calling
// Install and Uninstall to add and remove the correct shell installer commands
// retrospectively.
type Installer struct {
	installers []*Shell
	binaryPath string
}

// InstallerOptions represents a way to set optional values to a installer
// option.
// The InstallerOptions shows what options are available to change.
type InstallerOptions interface {
	AppendShell(*Shell)
	SetExecutable(Executable)
}

// InstallerOption captures a tweak that can be applied to the Installer.
type InstallerOption func(InstallerOptions)

type installer struct {
	shells     []*Shell
	executable Executable
}

func (i *installer) AppendShell(s *Shell) {
	i.shells = append(i.shells, s)
}

func (i *installer) SetExecutable(e Executable) {
	i.executable = e
}

func (i *installer) Executable() Executable {
	if i.executable == nil {
		return OSExecutable{}
	}
	return i.executable
}

// OptionShell allows the setting and appending of a shell option to configure
// the installer.
func OptionShell(shell *Shell) InstallerOption {
	return func(opt InstallerOptions) {
		opt.AppendShell(shell)
	}
}

// OptionExecutable allows the setting and appending of a exec option to
// configure the installer.
func OptionExecutable(exec Executable) InstallerOption {
	return func(opt InstallerOptions) {
		opt.SetExecutable(exec)
	}
}

// New creates a new installer with the correct dependencies and sane
// defaults.
// Returns an error if it can't locate the correct binary path.
func New(options ...InstallerOption) (*Installer, error) {
	opt := new(installer)
	for _, option := range options {
		option(opt)
	}

	path, err := opt.Executable().BinaryPath()
	if err != nil {
		return nil, err
	}

	return &Installer{
		installers: opt.shells,
		binaryPath: path,
	}, nil
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
