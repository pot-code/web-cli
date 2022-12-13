package task

import (
	"bytes"
	"fmt"
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
