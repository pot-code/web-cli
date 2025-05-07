package admingo

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/gomod"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type CreateModuleConfig struct {
	OutDir string `flag:"output" alias:"o" usage:"输出目录"`
	Name   string `arg:"0" alias:"MODULE_NAME" validate:"required,identifier"`
}

var CreateModuleCmd = command.NewCommand("module", "生成业务模块",
	&CreateModuleConfig{
		OutDir: "app",
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

		var buffers []*bytes.Buffer
		for range 4 {
			buffers = append(buffers, new(bytes.Buffer))
		}

		e := task.NewParallelScheduler()
		e.AddTask(
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/admingo/fx.go.tmpl"), buffers[0])).
				AddTask(task.NewTemplateRenderTask("fx", map[string]string{
					"projectName": projectName,
					"moduleName":  moduleName,
					"packageName": packageName,
				}, buffers[0], buffers[0])).
				AddTask(task.NewWriteFileToDiskTask("fx", ".go", outputDir, false, buffers[0])),
		)
		e.AddTask(
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/admingo/converter.go.tmpl"), buffers[1])).
				AddTask(task.NewTemplateRenderTask("converter", map[string]string{
					"packageName":   packageName,
					"outputPackage": outputPackage,
				}, buffers[1], buffers[1])).
				AddTask(task.NewWriteFileToDiskTask("converter", ".go", outputDir, false, buffers[1])),
		)
		e.AddTask(
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/admingo/schemas.go.tmpl"), buffers[2])).
				AddTask(task.NewTemplateRenderTask("schemas", map[string]string{
					"packageName": packageName,
				}, buffers[2], buffers[2])).
				AddTask(task.NewWriteFileToDiskTask("schemas", ".go", outputDir, false, buffers[2])),
		)
		e.AddTask(
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/admingo/service.go.tmpl"), buffers[3])).
				AddTask(task.NewTemplateRenderTask("service", map[string]string{
					"packageName": packageName,
				}, buffers[3], buffers[3])).
				AddTask(task.NewWriteFileToDiskTask("service", ".go", outputDir, false, buffers[3])),
		)
		return e.Run()
	}),
).Create()
