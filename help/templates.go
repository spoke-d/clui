package help

// HelpTemplateFormat defines a basic format for templating a help template.
const HelpTemplateFormat = `    {{green .Name}}	{{.Synopsis}}`

// HelpTemplate is the current view template of the help output.
const HelpTemplate = `
{{- if .Header }}
{{.Header}}
{{end}}
{{- if .Hint }}
Did you mean?
        {{green .Hint}}
{{end}}
Usage: {{green .Name}} [--version] [--help] [--debug] <command> [<args>]

{{- if gt (len .Commands) 0 }}

Available commands:
{{ range .Commands }}
%s
{{- end}}
{{- end}}

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`
