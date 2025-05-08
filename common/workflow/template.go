package workflow

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
)

type RenderEmbedTemplateWorkflow struct {
	templatePath string
	filename     string
	suffix       string
	data         map[string]string
	opts         []task.WriteFileOption
}

func NewRenderEmbedTemplateWorkflow(template string, data map[string]string, filename string, suffix string, opts ...task.WriteFileOption) *RenderEmbedTemplateWorkflow {
	return &RenderEmbedTemplateWorkflow{
		templatePath: template,
		data:         data,
		filename:     filename,
		suffix:       suffix,
		opts:         opts,
	}
}

func (r *RenderEmbedTemplateWorkflow) Run() error {
	b := new(bytes.Buffer)
	t := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider(r.templatePath), b)).
		AddTask(task.NewTemplateRenderTask("template", r.data, b, b)).
		AddTask(task.NewWriteFileToDiskTask(r.filename, r.suffix, b, r.opts...))
	return t.Run()
}
