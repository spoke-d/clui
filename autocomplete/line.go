package autocomplete

import "os"

// EnvComplete is used for the AutoComplete for getting the current terminal
// line.
const EnvComplete = "COMP_LINE"

// TerminalLine returns the current terminal line.
var TerminalLine = func() string {
	return os.Getenv(EnvComplete)
}
