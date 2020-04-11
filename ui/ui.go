package ui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/spoke-d/clui/ui/ask"
	"github.com/spoke-d/task/tomb"
)

// BasicUI is an implementation of UI that just outputs to the given writer.
type BasicUI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// NewBasicUI creates a new BasicUI with dependencies.
func NewBasicUI(stdin io.Reader, stdout, stderr io.Writer) *BasicUI {
	return &BasicUI{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// Ask asks the user for input using the given query. The response is
// returned as the given string, or an error.
func (u *BasicUI) Ask(query string) (string, error) {
	return u.ask(query, false)
}

// AskSecret asks the user for input using the given query, but does not echo
// the keystrokes to the terminal.
func (u *BasicUI) AskSecret(query string) (string, error) {
	return u.ask(query, true)
}

// Error is used for any error messages that might appear on standard
// error.
func (u *BasicUI) Error(message string) {
	fmt.Fprintln(u.stderr, message)
}

// Info is called for information related to the previous output.
// In general this may be the exact same as Output, but this gives
// UI implementors some flexibility with output formats.
func (u *BasicUI) Info(message string) {
	fmt.Fprintln(u.stdout, message)
}

// Output is called for normal standard output.
func (u *BasicUI) Output(template *Template, data interface{}) error {
	result, err := template.Render(data)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = fmt.Fprintln(u.stdout, strings.TrimSpace(result))
	return errors.WithStack(err)
}

func (u *BasicUI) ask(query string, secret bool) (string, error) {
	if _, err := fmt.Fprint(u.stdout, query+" "); err != nil {
		return "", err
	}

	// Register for interrupts so that we can catch it and immediately
	// return...
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	// Ask for input in a go-routine so that we can ignore it.
	lineCh := make(chan string, 1)

	t := tomb.New(false)
	err := t.Go(func(_ context.Context) error {
		var (
			line string
			err  error
		)
		if secret && isatty.IsTerminal(os.Stdin.Fd()) {
			asker := ask.New(u.stdin, u.stdout)
			line, err = asker.Password("")
		} else {
			r := bufio.NewReader(u.stdin)
			line, err = r.ReadString('\n')
		}
		if err != nil {
			return err
		}

		lineCh <- strings.TrimRight(line, "\r\n")
		return nil
	})
	if err != nil {
		return "", err
	}

	select {
	case line := <-lineCh:
		return line, nil
	case <-sigCh:
		// Print a newline so that any further output starts properly
		// on a new line.
		fmt.Fprintln(u.stdout)

		return "", errors.New("interrupted")
	}
}
