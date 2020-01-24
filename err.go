package clui

// Errno represents a error constants that can be reutrned from the CLI
type Errno int

// Code converts the Errno back into an int when type inference fails.
func (e Errno) Code() int {
	return int(e)
}

const (
	// EOK is non-standard representation of a success
	EOK Errno = 0

	// EPerm represents an operation not permitted
	EPerm Errno = 1

	// EKeyExpired is outside of POSIX 1, represents unknown error.
	EKeyExpired Errno = 127
)
