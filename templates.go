package clui

// TemplateVersion describes a template for rendering version.
const TemplateVersion = `
Client version: {{ .Version }}
`

// TemplateComplete describes a template for rendering autocomplete
// commands.
const TemplateComplete = `
{{- range $name := .}}
{{ $name }}
{{- end}}
`

// TemplatePlaceHolder describes a template for rendering placeholder text.
const TemplatePlaceHolder = `
This command is accessed by using one of the subcommands below.
`

// TemplateHelp represents a template for rendering the help for commands and
// subcommands.
const TemplateHelp = `
{{- if .Err }}
Found some issues:

    {{printf "%s\n" .Err | red}}
{{- end}}
{{- if .Hint }}
Did you mean?
    {{printf "%s\n" .Hint | green}}

{{- end}}
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

`
