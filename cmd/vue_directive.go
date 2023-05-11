package cmd

import (
	"bytes"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/file"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type VueDirectiveConfig struct {
	AddFolder bool   `flag:"add-folder" alias:"f" usage:"generate files in a folder with the name as the component"`
	OutDir    string `flag:"output" alias:"o" usage:"destination directory"`
	Name      string `arg:"0" alias:"COMPONENT_NAME" validate:"required,var"`
}

var VueDirectiveCmd = command.NewCliCommand("directive", "add vue directive",
	new(VueDirectiveConfig),
	command.WithArgUsage("DIRECTIVE_NAME"),
	command.WithAlias([]string{"d"}),
).AddHandlers(
	AddVueDirective,
).BuildCommand()

var AddVueDirective command.InlineHandler = func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*VueDirectiveConfig)
	fileName := config.Name
	camelName := strcase.ToCamel(fileName)
	outDir := config.OutDir

	b := new(bytes.Buffer)
	if err := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/vue/vue_directive.gotmpl"), b)).
		AddTask(task.NewTemplateRenderTask("vue_directive", map[string]string{
			"name": camelName,
		}, b, b)).
		AddTask(task.NewWriteFileToDiskTask(fileName, file.TypescriptSuffix, outDir, false, b)).
		Run(); err != nil {
		return err
	}
	return nil
}
