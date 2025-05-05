package react

import (
	"bytes"
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type ReactHookConfig struct {
	Name    string `arg:"0" alias:"HOOK_NAME" validate:"required,var"`
	OutDir  string `flag:"output" alias:"o" usage:"destination directory"`
	AddTest bool   `flag:"add-test" alias:"t" usage:"add associated hook test file"`
}

var ReactHookCmd = command.NewCommand("hook", "add react hook",
	new(ReactHookConfig),
	command.WithAlias([]string{"k"}),
).AddHandler(
	AddReactHook,
).Create()

var AddReactHook = command.InlineHandler[*ReactHookConfig](func(c *cli.Context, config *ReactHookConfig) error {
	varName := strcase.ToCamel(config.Name)
	filename := strcase.ToKebab(fmt.Sprintf("use%s", varName))

	b1 := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook.go.tmpl"), b1)).
			AddTask(task.NewTemplateRenderTask("react_hook", map[string]string{"name": varName}, b1, b1)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".ts", config.OutDir, false, b1)),
	}

	if config.AddTest {
		b2 := new(bytes.Buffer)
		tasks = append(tasks,
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook_test.go.tmpl"), b2)).
				AddTask(task.NewTemplateRenderTask("react_hook_test", map[string]string{"name": varName}, b2, b2)).
				AddTask(task.NewWriteFileToDiskTask(filename, ".test.ts", config.OutDir, false, b2)))
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})
