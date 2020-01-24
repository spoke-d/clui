package ui

import (
	"bufio"
	"fmt"
	"github.com/spoke-d/clui/ui/style"
	"strings"
)

func indent(input string) string {
	return indentTimes(input, 1)
}

func indentTimes(input string, count int) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)

	indent := strings.Repeat("    ", count)

	var result string
	for scanner.Scan() {
		result += fmt.Sprintf("%s%s\n", indent, scanner.Text())
	}
	return result
}

var (
	styleRed   = style.New(style.FgRed)
	styleGreen = style.New(style.FgGreen)
)

func red(color bool) func(string) string {
	return func(input string) string {
		if !color {
			return input
		}
		return styleRed.Format(input)
	}
}

func green(color bool) func(string) string {
	return func(input string) string {
		if !color {
			return input
		}
		return styleGreen.Format(input)
	}
}
