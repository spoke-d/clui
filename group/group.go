package group

import (
	"fmt"
	"math"
	"strings"

	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/group/distance"
	"github.com/spoke-d/clui/radix"
	task "github.com/spoke-d/task/group"
)

// PlaceHolder is a type which is called when constructing the group and a place
// holder command is required.
type PlaceHolder func(string) Command

// Command represents an abstraction of command.
type Command interface {

	// Flags returns the FlagSet associated with the command. All the flags are
	// parsed before running the command.
	FlagSet() *flagset.FlagSet

	// Help should return a long-form help text that includes the command-line
	// usage. A brief few sentences explaining the function of the command, and
	// the complete list of flags the command accepts.
	Help() string

	// Synopsis should return a one-line, short synopsis of the command.
	// This should be short (50 characters of less ideally).
	Synopsis() string

	// Init is called with all the args required to run a command.
	// This is separated from Run, to allow the preperation of a command, before
	// it's run.
	Init([]string, bool) error

	// Run should run the actual command with the given CLI instance and
	// command-line arguments. It should return the exit status when it is
	// finished.
	//
	// There are a handful of special exit codes that can return documented
	// behavioral changes.
	Run(*task.Group)
}

// GroupOptions represents a way to set optional values to a autocomplete
// option.
// The GroupOptions shows what options are available to change.
type GroupOptions interface {
	SetPlaceHolder(PlaceHolder)
}

// GroupOption captures a tweak that can be applied to the Group.
type GroupOption func(GroupOptions)

type grp struct {
	placeHolder PlaceHolder
}

func (s *grp) SetPlaceHolder(p PlaceHolder) {
	s.placeHolder = p
}

// OptionPlaceHolder allows the setting a place holder option to configure
// the group.
func OptionPlaceHolder(i PlaceHolder) GroupOption {
	return func(opt GroupOptions) {
		opt.SetPlaceHolder(i)
	}
}

// Group holds the commands in a central repository for easy access.
type Group struct {
	commands      map[string]Command
	commandTree   *radix.Tree
	placeholderFn PlaceHolder
}

// New creates a Group with sane defaults.
func New(options ...GroupOption) *Group {
	opt := new(grp)
	for _, option := range options {
		option(opt)
	}

	return &Group{
		commands:      make(map[string]Command),
		commandTree:   radix.New(),
		placeholderFn: opt.placeHolder,
	}
}

// Add a Command to the Group for a given key. The key is normalized
// to remove trailing spaces for consistency.
// Returns an error when inserting into the Group fails
func (r *Group) Add(key string, cmd Command) error {
	k := normalizeKey(key)

	if _, _, err := r.commandTree.Insert(k, cmd); err != nil {
		return err
	}
	r.commands[k] = cmd
	return nil
}

// Remove a Command from the Group for a given key. The key is
// normalized to remove trailing spaces for consistency.
// Returns an error when deleting from the Group upon failure
func (r *Group) Remove(key string) (Command, error) {
	k := normalizeKey(key)

	cmd, ok := r.commands[k]
	if ok {
		delete(r.commands, k)
	}

	if _, v := r.commandTree.Delete(k); ok && v {
		return cmd, nil
	}

	return nil, fmt.Errorf("no valid key found for %q", key)
}

// Get returns a Command for a given key. The key is normalized to remove
// trailing spaces for consistency.
// Returns true if it was found.
func (r *Group) Get(key string) (Command, bool) {
	k := normalizeKey(key)
	if _, ok := r.commands[k]; !ok {
		return nil, false
	}
	raw, ok := r.commandTree.Get(k)
	if !ok {
		return nil, false
	}
	cmd, ok := raw.(Command)
	if !ok {
		return nil, false
	}
	return cmd, true
}

// GetClosestName returns the closest command to the given key
func (r *Group) GetClosestName(key string) (string, bool) {
	if len(key) == 0 {
		return "", false
	}
	closest := struct {
		name     string
		distance int
	}{
		distance: math.MaxInt64,
	}
	r.commandTree.Walk(func(name string, value radix.Value) bool {
		d := distance.ComputeDistance(key, name)
		if strings.HasPrefix(name, key[:1]) && d < closest.distance {
			closest.name = name
			closest.distance = d
		}
		return false
	})
	if closest.name == "" {
		return "", false
	}
	k := normalizeKey(closest.name)
	if _, ok := r.commands[k]; !ok {
		return "", false
	}
	return k, true
}

// WalkPrefix is used to walk the tree under a prefix
func (r *Group) WalkPrefix(prefix string, fn radix.WalkFn) {
	r.commandTree.WalkPrefix(prefix, fn)
}

// LongestPrefix is like Get, but instead of an exact match, it will return
// the longest prefix match.
func (r *Group) LongestPrefix(key string) (string, bool) {
	k := normalizeKey(key)
	s, _, ok := r.commandTree.LongestPrefix(k)
	return s, ok
}

// Process runs through the registry and fills in any commands that are required
// for nesting (sub commands)
// Returns an error if there was an issue adding any commands to the underlying
// storage.
func (r *Group) Process() error {
	if r.Nested() {
		var (
			walkFn radix.WalkFn
			insert = make(map[string]struct{})
		)

		walkFn = func(k string, v radix.Value) bool {
			idx := strings.LastIndex(k, " ")
			if idx == -1 {
				// If there is no space, just ignore top level commands
				return false
			}

			// Trim up to that space so we can get the expected parent
			k = k[:idx]
			if _, ok := r.Get(k); ok {
				return false
			}

			// We're missing the parent, so let's insert this
			insert[k] = struct{}{}

			// Call the walk function recursively so we check this one too
			return walkFn(k, nil)
		}
		r.commandTree.Walk(walkFn)

		// Insert any that we're missing
		for k := range insert {
			if err := r.Add(k, r.placeholderFn(k)); err != nil {
				return err
			}
		}
	}
	return nil
}

// Nested returns if the commands with in the group are nested in anyway.
func (r *Group) Nested() bool {
	for k := range r.commands {
		if strings.ContainsRune(k, ' ') {
			return true
		}
	}
	return false
}

// normalizeKey attempts to normalize a command key before inserting it into a
// set of maps/trees. This should help with any possible inconsistencies when
// querying the structure.
func normalizeKey(k string) string {
	s := strings.Split(strings.TrimSpace(k), " ")

	var names []string
	for _, v := range s {
		t := strings.TrimSpace(v)
		if t != "" {
			names = append(names, t)
		}
	}
	return strings.Join(names, " ")
}
