package main

import (
	"flag"

	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/commands"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/ui"
	"github.com/spoke-d/task/group"
)

type configShowCmd struct {
	ui      clui.UI
	flagSet *flagset.FlagSet

	template string
	server   bool
}

func configShowCmdFn(ui clui.UI) clui.Command {
	cmd := &configShowCmd{
		ui:      ui,
		flagSet: flagset.New("config show", flag.ContinueOnError),
	}
	cmd.init()
	return cmd
}

func (v *configShowCmd) init() {
	v.flagSet.StringVar(&v.template, "template", "{{.Key}}	{{.Value}}", "Template for show key and values")
}

func (v *configShowCmd) FlagSet() *flagset.FlagSet {
	return v.flagSet
}

func (v *configShowCmd) Help() string {
	return `
Show the current key and values of a configuration found
on the server.
`
}

func (v *configShowCmd) Synopsis() string {
	return "Show configuration."
}

type config struct {
	Key   string
	Value interface{}
}

func (v *configShowCmd) Run(g *group.Group) {
	template := ui.NewTemplate(configShowTemplate, ui.OptionFormat(v.template))
	v.ui.Output(template, struct {
		Config []config
	}{
		Config: []config{
			{
				Key:   "default",
				Value: "(null)",
			},
		},
	})
	commands.Nothing(g)
}

const configShowTemplate = `
Configuration: 

{{range .Config -}}
%s
{{- end}}
`
