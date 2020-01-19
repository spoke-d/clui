package install

import "os/user"

// OSUser represents a current operating system user.
type OSUser struct {
	user *user.User
}

// CurrentUser returns the current operating system user.
// If there is an error trying to get the current user, it returns an error.
func CurrentUser() (OSUser, error) {
	u, err := user.Current()
	if err != nil {
		return OSUser{}, err
	}

	return OSUser{
		user: u,
	}, nil
}

// HomeDir returns the current home directory for the os user.
func (u OSUser) HomeDir() string {
	return u.user.HomeDir
}
