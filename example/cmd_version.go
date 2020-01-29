package main

import (
	"flag"

	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/commands"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/ui"
	"github.com/spoke-d/task/group"
)

type versionCmd struct {
	ui      clui.UI
	flagSet *flagset.FlagSet

	template string
	server   bool
}

func versionCmdFn(ui clui.UI) clui.Command {
	cmd := &versionCmd{
		ui:      ui,
		flagSet: flagset.New("version", flag.ContinueOnError),
	}
	cmd.init()
	return cmd
}

func (v *versionCmd) init() {
	v.flagSet.StringVar(&v.template, "template", "{{.Version}}", "Template for the version template")
	v.flagSet.BoolVar(&v.server, "server", false, "Request server version")
}

func (v *versionCmd) FlagSet() *flagset.FlagSet {
	return v.flagSet
}

func (v *versionCmd) Help() string {
	return `
Show the current client version along with optionally
requesting the server version.

The two version can be safely out of sync, as long as
the API versions can correctly interact with the said
version.
`
}

func (v *versionCmd) Synopsis() string {
	return "Show client and server version."
}

func (v *versionCmd) Init([]string, bool) error {
	return nil
}

func (v *versionCmd) Run(g *group.Group) {
	template := ui.NewTemplate(versionTemplate, ui.OptionFormat(v.template))
	v.ui.Output(template, struct {
		Version string
	}{
		Version: "1.0.0",
	})
	commands.Nothing(g)
}

const versionTemplate = `
Version: %s
`
