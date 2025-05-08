package admingo

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pot-code/web-cli/common/workflow"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/gomod"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type CreateModuleConfig struct {
	OutDir string `flag:"output" alias:"o" usage:"输出目录"`
	Name   string `arg:"0" alias:"MODULE_NAME" validate:"required,identifier"`
}

var CreateModuleCmd = command.NewCommand("module", "生成业务模块",
	&CreateModuleConfig{
		OutDir: "internal/app",
	},
	command.WithAlias([]string{"m"}),
).AddHandler(
	command.InlineHandler[*CreateModuleConfig](func(c *cli.Context, config *CreateModuleConfig) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current working directory: %w", err)
		}

		gm, err := gomod.ParseGoMod(path.Join(cwd, "go.mod"))
		if err != nil {
			return err
		}

		projectName := gm.ProjectName()
		packageName := strings.ToLower(config.Name)
		moduleName := packageName + "s"
		outputDir := path.Join(cwd, config.OutDir, packageName)
		outputPackage := path.Join(projectName, config.OutDir, packageName)

		e := task.NewParallelScheduler()
		e.AddTask(
			workflow.NewRenderEmbedTemplateWorkflow("templates/admingo/fx.go.tmpl", map[string]string{
				"projectName": projectName,
				"moduleName":  moduleName,
				"packageName": packageName,
			}, "fx", ".go", task.WithFolder(outputDir)),
		)
		e.AddTask(
			workflow.NewRenderEmbedTemplateWorkflow("templates/admingo/converter.go.tmpl", map[string]string{
				"packageName":   packageName,
				"outputPackage": outputPackage,
			}, "converter", ".go", task.WithFolder(outputDir)),
		)
		e.AddTask(
			workflow.NewRenderEmbedTemplateWorkflow("templates/admingo/schemas.go.tmpl", map[string]string{
				"packageName": packageName,
			}, "schemas", ".go", task.WithFolder(outputDir)),
		)
		e.AddTask(
			workflow.NewRenderEmbedTemplateWorkflow("templates/admingo/service.go.tmpl", map[string]string{
				"packageName": packageName,
			}, "service", ".go", task.WithFolder(outputDir)),
		)
		return e.Run()
	}),
).Create()
