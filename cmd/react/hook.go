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

var ReactHookCmd = command.NewCommandBuilder("hook", "add react hook",
	new(ReactHookConfig),
	command.WithArgUsage("HOOK_NAME"),
	command.WithAlias([]string{"k"}),
).AddHandlers(
	AddReactHook,
).Build()

var AddReactHook = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	rhc := cfg.(*ReactHookConfig)
	varName := strcase.ToCamel(rhc.Name)
	filename := strcase.ToLowerCamel(fmt.Sprintf("use%s", varName))

	b1 := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook.gotmpl"), b1)).
			AddTask(task.NewTemplateRenderTask("react_hook", map[string]string{"name": varName}, b1, b1)).
			AddTask(task.NewWriteFileToDiskTask(filename, ".ts", rhc.OutDir, false, b1)),
	}

	if rhc.AddTest {
		b2 := new(bytes.Buffer)
		tasks = append(tasks,
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook_test.gotmpl"), b2)).
				AddTask(task.NewTemplateRenderTask("react_hook_test", map[string]string{"name": varName}, b2, b2)).
				AddTask(task.NewWriteFileToDiskTask(filename, ".test.ts", rhc.OutDir, false, b2)))
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})
