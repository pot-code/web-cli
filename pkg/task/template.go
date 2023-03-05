package task

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type GenerateFileFromTemplateTask struct {
	ft *WriteFileToDiskTask
	tr *TemplateRenderTask
}

func NewGenerateFileFromTemplateTask(
	fileName string,
	suffix string,
	folder string,
	overwrite bool,
	templateName string,
	templateProvider TemplateProvider,
	templateData interface{}) *GenerateFileFromTemplateTask {
	b := new(bytes.Buffer)
	return &GenerateFileFromTemplateTask{
		NewWriteFileToDiskTask(fileName, suffix, folder, overwrite, b),
		NewTemplateRenderTask(templateName, templateProvider, templateData, b),
	}
}

func (t *GenerateFileFromTemplateTask) Run() error {
	err := NewSequentialScheduler().
		AddTask(t.tr).
		AddTask(t.ft).
		Run()
	if err != nil {
		return fmt.Errorf("run GenerateFileFromTemplateTask: %w", err)
	}
	return nil
}

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
		return fmt.Errorf("get template from provider: %w", err)
	}

	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return fmt.Errorf("read template data: %w", err)
	}

	log.WithFields(log.Fields{"template_name": trt.Name, "data": trt.Data}).Debug("render template")
	err = RenderTextTemplate(&RenderRequest{
		Name:     trt.Name,
		Template: string(b),
		Data:     trt.Data,
	}, trt.out)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	return nil
}

type RenderRequest struct {
	// template name
	Name     string
	Template string
	Data     interface{}
}

func RenderTextTemplate(req *RenderRequest, out io.Writer) error {
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
