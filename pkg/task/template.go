package task

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pot-code/web-cli/pkg/template"
	log "github.com/sirupsen/logrus"
)

type TemplateProvider interface {
	Get() (io.Reader, error)
}

var _ Task = (*TemplateRenderTask)(nil)

type TemplateRenderTask struct {
	// template name
	Name     string
	Provider TemplateProvider
	Data     interface{}
	out      io.Writer
}

func NewTemplateRenderTask(name string, provider TemplateProvider, data interface{}, out io.Writer) *TemplateRenderTask {
	return &TemplateRenderTask{
		Name:     name,
		Provider: provider,
		Data:     data,
		out:      out,
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
	if trt.Name == "" {
		return errors.New("empty template name")
	}
	return nil
}

func (trt *TemplateRenderTask) renderTemplate() error {
	p := trt.Provider

	fd, err := p.Get()
	if err != nil {
		return fmt.Errorf("get template: %w", err)
	}

	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return fmt.Errorf("read template data: %w", err)
	}

	log.WithFields(log.Fields{"template_name": trt.Name, "data": trt.Data}).Debug("render template")
	err = template.RenderTextTemplate(&template.RenderRequest{
		Name:     trt.Name,
		Template: string(b),
		Data:     trt.Data,
	}, trt.out)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	return nil
}
