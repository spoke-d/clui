package ui

import (
	"bytes"
	"fmt"
	"io"
	ttemplate "text/template"

	"github.com/pkg/errors"
)

// TemplateOptions represents a way to set optional values to a template
// option.
// The TemplateOptions shows what options are available to change.
type TemplateOptions interface {
	SetName(string)
	SetFormat(string)
	SetColor(bool)
}

// TemplateOption captures a tweak that can be applied to the Template.
type TemplateOption func(TemplateOptions)

type template struct {
	name   string
	format string
	color  bool
}

func (s *template) SetName(p string) {
	s.name = p
}

func (s *template) SetFormat(p string) {
	s.format = p
}

func (s *template) SetColor(p bool) {
	s.color = p
}

func (s *template) Name() string {
	if s.name == "" {
		return "view"
	}
	return s.name
}

// OptionName allows the setting a name option to configure the template.
func OptionName(i string) TemplateOption {
	return func(opt TemplateOptions) {
		opt.SetName(i)
	}
}

// OptionFormat allows the setting a format option to configure the template.
func OptionFormat(i string) TemplateOption {
	return func(opt TemplateOptions) {
		opt.SetFormat(i)
	}
}

// OptionColor allows the setting a color option to configure the template.
func OptionColor(i bool) TemplateOption {
	return func(opt TemplateOptions) {
		opt.SetColor(i)
	}
}

// Template represents a view that will be rendered by the UI.
type Template struct {
	format   string
	template string
	renderer *ttemplate.Template
}

// NewTemplate creates a template for rendering a view.
func NewTemplate(t string, options ...TemplateOption) *Template {
	opt := new(template)
	for _, option := range options {
		option(opt)
	}

	renderer := ttemplate.New(opt.Name())
	renderer.Funcs(map[string]interface{}{
		"indent": indent,
		"red":    red(opt.color),
		"green":  green(opt.color),
	})

	return &Template{
		format:   opt.format,
		template: t,
		renderer: renderer,
	}
}

// Render renders the template with given data.
func (t *Template) Render(data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.Write(&buf, data); err != nil {
		return "", errors.WithStack(err)
	}
	return buf.String(), nil
}

func (t *Template) Write(writer io.Writer, data interface{}) error {
	view := t.template
	if t.format != "" {
		view = fmt.Sprintf(t.template, t.format)
	}

	template, err := t.renderer.Parse(view)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := template.Execute(writer, data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
