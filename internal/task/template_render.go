package task

import (
	"fmt"
	"io"
	"text/template"
)

type RenderRequest struct {
	// template name
	Name     string
	Template string
	Data     interface{}
}

type DefaultTemplateRenderer struct{}

func NewDefaultTemplateRenderer() *DefaultTemplateRenderer {
	return &DefaultTemplateRenderer{}
}

func (r *DefaultTemplateRenderer) Render(req *RenderRequest, out io.Writer) error {
	t := template.New(req.Name)

	pt, err := t.Parse(req.Template)
	if err != nil {
		return fmt.Errorf("parse template [name: %s]: %w", req.Name, err)
	}

	if err := pt.Execute(out, req.Data); err != nil {
		return fmt.Errorf("execute template [name: %s]: %w", req.Name, err)
	}
	return nil
}
