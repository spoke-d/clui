package clui

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spoke-d/clui/group"
	"github.com/spoke-d/clui/radix"
)

// FindChildren returns the sub commands.
// This will only contain immediate sub commands.
func FindChildren(commands *group.Group, prefix string, includeSubKeys bool) (map[string]Command, error) {
	// if our prefix isn't empty, make sure it ends in ' '
	if prefix != "" && prefix[len(prefix)-1] != ' ' {
		prefix += " "
	}

	// Get all the subkeys of this command
	var keys []string
	commands.WalkPrefix(prefix, func(k string, v radix.Value) bool {
		// Ignore any sub-sub keys, i.e. "foo bar baz" when we want "foo bar"
		if !includeSubKeys && strings.Contains(k[len(prefix):], " ") {
			return false
		}

		keys = append(keys, k)

		return false
	})

	// For each of the keys return that in the map
	res := make(map[string]Command, len(keys))
	for _, k := range keys {
		cmd, ok := commands.Get(k)
		if !ok {
			return nil, errors.Errorf("not found: %q", k)
		}
		res[k] = cmd
	}

	return res, nil
}
