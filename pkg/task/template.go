package task

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	log "github.com/sirupsen/logrus"
)

// type GenerateFromTemplateTask struct {
// 	ft *WriteFileToDiskTask
// 	tr *TemplateRenderTask
// }

// func NewGenerateFromTemplateTask(
// 	fileName string,
// 	suffix string,
// 	folder string,
// 	overwrite bool,
// 	templateName string,
// 	templateProvider TemplateProvider,
// 	templateData interface{}) *GenerateFromTemplateTask {
// 	b := new(bytes.Buffer)
// 	return &GenerateFromTemplateTask{
// 		NewWriteFileToDiskTask(fileName, suffix, folder, overwrite, b),
// 		NewTemplateRenderTask(templateName, templateProvider, templateData, b),
// 	}
// }

// func (t *GenerateFromTemplateTask) Run() error {
// 	err := NewSequentialScheduler().
// 		AddTask(t.tr).
// 		AddTask(t.ft).
// 		Run()
// 	if err != nil {
// 		return fmt.Errorf("run GenerateFileFromTemplateTask: %w", err)
// 	}
// 	return nil
// }

// type TemplateProvider interface {
// 	Get() (io.Reader, error)
// }

type TemplateRenderTask struct {
	// template name
	name string
	data interface{}
	in   io.Reader
	out  io.Writer
}

func NewTemplateRenderTask(name string, data interface{}, in io.Reader, out io.Writer) *TemplateRenderTask {
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
	b, err := ioutil.ReadAll(trt.in)
	if err != nil {
		return fmt.Errorf("read template data: %w", err)
	}

	log.WithFields(log.Fields{"template_name": trt.name, "data": trt.data}).Debug("render template")
	err = RenderTextTemplate(&RenderRequest{
		Name:     trt.name,
		Template: string(b),
		Data:     trt.data,
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
