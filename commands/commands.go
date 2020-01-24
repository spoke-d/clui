package commands

import "github.com/spoke-d/task/group"

// Disguard is used as a High Order Function to disguard error messages
var Disguard = func(err error) {}

// Nothing is used as a High Order Function to consume all group event runs into
// nothing.
var Nothing = func(g *group.Group) {
	g.Add(func() error { return nil }, Disguard)
}
