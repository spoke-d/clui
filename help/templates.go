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
{{- if .ShowHelp }}
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
{{- end}}
`

// CommandHelpTemplate represents a template for rendering the help for commands
// and subcommands.
const CommandHelpTemplate = `
{{- if .Err }}
Found some issues:

    {{red .Err}}

See {{ print " --help" | print .Name | green }} for more information.
{{end -}}
{{- if .Hint }}
Did you mean?
    {{printf "%s\n" .Hint | green}}

{{end -}}
{{- if .ShowHelp }}
{{- if .Name}}
Usage:

    {{green .Name}}{{- if .Flags}} [flags]{{- end}}

{{- if gt (len .Flags) 0 }}
{{range $flag := .Flags }}
    {{green $.Name}} {{$flag}}
{{- end}}
{{- end}}
{{- if gt (len .Usages) 0 }}
{{range $usage := .Usages }}
    {{green $.Name}} {{$usage}}
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
{{- end}}
`
