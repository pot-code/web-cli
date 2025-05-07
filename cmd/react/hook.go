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
	Name    string `arg:"0" alias:"HOOK_NAME" validate:"required,identifier"`
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

	var buffers []*bytes.Buffer
	for range 2 {
		buffers = append(buffers, new(bytes.Buffer))
	}

	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook.go.tmpl"), buffers[0])).
			AddTask(task.NewTemplateRenderTask("react_hook", map[string]string{"name": varName}, buffers[0], buffers[0])).
			AddTask(task.NewWriteFileToDiskTask(filename, ".ts", buffers[0], task.WithFolder(config.OutDir))),
	}

	if config.AddTest {
		tasks = append(tasks,
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook_test.go.tmpl"), buffers[1])).
				AddTask(task.NewTemplateRenderTask("react_hook_test", map[string]string{"name": varName}, buffers[1], buffers[1])).
				AddTask(task.NewWriteFileToDiskTask(filename, ".test.ts", buffers[1], task.WithFolder(config.OutDir))))
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})
