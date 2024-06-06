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

type VueDirectiveConfig struct {
	Isolated bool   `flag:"isolated" alias:"i" usage:"generate files in a folder"`
	OutDir   string `flag:"output" alias:"o" usage:"destination directory"`
	Name     string `arg:"0" alias:"COMPONENT_NAME" validate:"required,var"`
}

var VueDirectiveCmd = command.NewBuilder("directive", "add vue directive",
	new(VueDirectiveConfig),
	command.WithArgUsage("DIRECTIVE_NAME"),
	command.WithAlias([]string{"d"}),
).AddHandlers(
	AddVueDirective,
).Build()

var AddVueDirective = command.InlineHandler[*VueDirectiveConfig](func(c *cli.Context, config *VueDirectiveConfig) error {
	filename := config.Name
	camelName := strcase.ToCamel(filename)
	outDir := config.OutDir

	if config.Isolated {
		outDir = path.Join(config.OutDir, filename)
		filename = "index"
	}

	b := new(bytes.Buffer)
	if err := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/vue/vue_directive.gotmpl"), b)).
		AddTask(task.NewTemplateRenderTask("vue_directive", map[string]string{
			"name": camelName,
		}, b, b)).
		AddTask(task.NewWriteFileToDiskTask(filename, ".ts", outDir, false, b)).
		Run(); err != nil {
		return err
	}
	return nil
})
