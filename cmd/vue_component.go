package cmd

import (
	"bytes"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/file"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type VueComponentConfig struct {
	AddFolder bool   `flag:"add-folder" alias:"f" usage:"generate files in a folder with the name as the component"`
	OutDir    string `flag:"output" alias:"o" usage:"destination directory"`
	Name      string `arg:"0" alias:"COMPONENT_NAME" validate:"required,var"`
}

var VueComponentCmd = command.NewCliCommand("component", "add vue component",
	new(VueComponentConfig),
	command.WithArgUsage("COMPONENT_NAME"),
	command.WithAlias([]string{"c"}),
).AddHandlers(
	AddVueComponent,
).BuildCommand()

var AddVueComponent command.InlineHandler = func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*VueComponentConfig)
	fileName := strcase.ToCamel(config.Name)
	kebabName := strcase.ToKebab(config.Name)
	outDir := config.OutDir

	if config.AddFolder {
		fileName = "index"
		outDir = path.Join(config.OutDir, kebabName)
	}

	b := new(bytes.Buffer)
	if err := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/vue/vue_component.gotmpl"), b)).
		AddTask(task.NewTemplateRenderTask("vue_component", nil, b, b)).
		AddTask(task.NewWriteFileToDiskTask(fileName, file.VueComponentSuffix, outDir, false, b)).
		Run(); err != nil {
		return err
	}
	return nil
}
