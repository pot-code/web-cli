package task

import (
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/rs/zerolog/log"
)

type TemplateRenderTask struct {
	// name template name
	name string
	// data template data
	data any
	in   io.Reader
	out  io.Writer
}

func NewTemplateRenderTask(name string, data any, in io.Reader, out io.Writer) *TemplateRenderTask {
	return &TemplateRenderTask{
		name: name,
		data: data,
		in:   in,
		out:  out,
	}
}

func (trt *TemplateRenderTask) Run() error {
	if err := trt.validateTask(); err != nil {
		return err
	}

	if err := trt.renderTemplate(); err != nil {
		return err
	}
	return nil
}

func (trt *TemplateRenderTask) validateTask() error {
	if trt.name == "" {
		return errors.New("empty template name")
	}
	return nil
}

func (trt *TemplateRenderTask) renderTemplate() error {
	b, err := io.ReadAll(trt.in)
	if err != nil {
		return fmt.Errorf("read template data: %w", err)
	}

	log.Debug().Str("task", "TemplateRenderTask").Str("template_name", trt.name).Interface("data", trt.data).Msg("render template")
	err = RenderTextTemplate(&RenderRequest{
		name:     trt.name,
		template: string(b),
		data:     trt.data,
	}, trt.out)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	return nil
}

type RenderRequest struct {
	// name template name
	name     string
	template string
	data     any
}

func RenderTextTemplate(req *RenderRequest, out io.Writer) error {
	t := template.New(req.name)
	pt, err := t.Parse(req.template)
	if err != nil {
		return fmt.Errorf("parse template [name: %s]: %w", req.name, err)
	}
	if err := pt.Execute(out, req.data); err != nil {
		return fmt.Errorf("execute template [name: %s]: %w", req.name, err)
	}
	return nil
}
