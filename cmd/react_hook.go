package cmd

import (
	"bytes"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

const TypescriptSuffix = ".ts"

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
	config := cfg.(*ReactHookConfig)
	name := strcase.ToLowerCamel(config.Name)

	b1 := new(bytes.Buffer)
	tasks := []task.Task{
		task.NewSequentialScheduler().
			AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react_hook.gohtml"), b1)).
			AddTask(task.NewTemplateRenderTask("react_hook", map[string]string{"name": name}, b1, b1)).
			AddTask(task.NewWriteFileToDiskTask(name, TypescriptSuffix, config.OutDir, false, b1)),
	}

	if config.AddTest {
		b2 := new(bytes.Buffer)
		tasks = append(tasks,
			task.NewSequentialScheduler().
				AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/react_hook_test.gohtml"), b2)).
				AddTask(task.NewTemplateRenderTask("react_hook_test", map[string]string{"name": name}, b2, b2)).
				AddTask(task.NewWriteFileToDiskTask(name, ReactTestSuffix, config.OutDir, false, b2)))
	}

	s := task.NewParallelScheduler()
	for _, t := range tasks {
		s.AddTask(t)
	}
	return s.Run()
})
