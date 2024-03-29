package cmd

import (
	"bytes"
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/file"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type ReactHookConfig struct {
	Name    string `arg:"0" alias:"HOOK_NAME" validate:"required,var"`
	OutDir  string `flag:"output" alias:"o" usage:"destination directory"`
	AddTest bool   `flag:"add-test" alias:"t" usage:"add associated hook test file"`
}

var ReactHookCmd = command.NewCliCommand("hook", "add react hook",
	new(ReactHookConfig),
	command.WithArgUsage("HOOK_NAME"),
	command.WithAlias([]string{"k"}),
).AddHandlers(
	AddReactHook,
).BuildCommand()

var AddReactHook = command.InlineHandler(func(c *cli.Context, cfg interface{}) error {
	rhc := cfg.(*ReactHookConfig)
	varName := strcase.ToCamel(rhc.Name)
	fileName := strcase.ToKebab(fmt.Sprintf("use%s", varName))

	b1 := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook.gotmpl"), b1)).
			AddTask(task.NewTemplateRenderTask("react_hook", map[string]string{"name": varName}, b1, b1)).
			AddTask(task.NewWriteFileToDiskTask(fileName, file.TypescriptSuffix, rhc.OutDir, false, b1)),
	}

	if rhc.AddTest {
		b2 := new(bytes.Buffer)
		tasks = append(tasks,
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react/react_hook_test.gotmpl"), b2)).
				AddTask(task.NewTemplateRenderTask("react_hook_test", map[string]string{"name": varName}, b2, b2)).
				AddTask(task.NewWriteFileToDiskTask(fileName, file.TypescriptTestSuffix, rhc.OutDir, false, b2)))
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})
