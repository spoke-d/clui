package help

// HelpTemplateFormat defines a basic format for templating a help template.
const HelpTemplateFormat = `    {{green .Name}}	{{.Synopsis}}`

// BasicHelpTemplate is the current view template of the help output.
const BasicHelpTemplate = `
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

// CommandHelpTemplate represents a template for rendering the help for commands
// and subcommands.
const CommandHelpTemplate = `
{{- if .Err }}
Found some issues:

    {{red .Err}}
{{end -}}
{{- if .Hint }}
Did you mean?
    {{printf "%s\n" .Hint | green}}

{{end -}}
{{- if .Name}}
Usage:

    {{green .Name}}{{- if .Flags}} [flags]{{- end}}

{{- if gt (len .Flags) 0 }}
{{range $flag := .Flags }}
    {{green $.Name}}  {{$flag}}
{{- end}}
{{- end}}
{{- end}}

Description:
    {{ indent .Help }}

{{- if gt (len .Commands) 0 }}

Available commands:
{{ range .Commands }}
    {{green .Name}}	{{.Synopsis}}
{{- end}}
{{- end}}

Global Flags:

        --debug        Show all debug messages
    -h, --help         Print command help
        --version      Print client version
`
