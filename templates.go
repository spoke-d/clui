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
Selection of commands that are nested under this one.
`

// TemplateFlags describes a template for rendering flags in help.
const TemplateFlags = `
{{.Name}}	{{.Usage}} (defaults: "{{.Defaults}}")
`
