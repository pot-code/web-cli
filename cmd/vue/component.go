package vue

import (
	"bytes"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type VueComponentConfig struct {
	Isolated bool   `flag:"isolated" alias:"i" usage:"generate files in a folder"`
	OutDir   string `flag:"output" alias:"o" usage:"destination directory"`
	Name     string `arg:"0" alias:"COMPONENT_NAME" validate:"required,var"`
}

var VueComponentCmd = command.NewBuilder("component", "add vue component",
	new(VueComponentConfig),
	command.WithArgUsage("COMPONENT_NAME"),
	command.WithAlias([]string{"c"}),
).AddHandlers(
	AddVueComponent,
).Build()

var AddVueComponent command.InlineHandler = func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*VueComponentConfig)
	filename := strcase.ToCamel(config.Name)
	kebabName := strcase.ToKebab(config.Name)
	outDir := config.OutDir

	if config.Isolated {
		filename = "index"
		outDir = path.Join(config.OutDir, kebabName)
	}

	b := new(bytes.Buffer)
	if err := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/vue/vue_component.gotmpl"), b)).
		AddTask(task.NewTemplateRenderTask("vue_component", nil, b, b)).
		AddTask(task.NewWriteFileToDiskTask(filename, ".vue", outDir, false, b)).
		Run(); err != nil {
		return err
	}
	return nil
}
