package vue

import (
	"bytes"
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/internal/command"
	"github.com/pot-code/web-cli/internal/provider"
	"github.com/pot-code/web-cli/internal/task"
	"github.com/urfave/cli/v2"
)

type VueUseStoreConfig struct {
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
	Name   string `arg:"0" alias:"MODULE_NAME" validate:"required,var"`
}

var VueUseStoreCmd = command.NewCommandBuilder("store", "add vue pinia store",
	new(VueUseStoreConfig),
	command.WithArgUsage("MODULE_NAME"),
	command.WithAlias([]string{"s"}),
).AddHandlers(
	UseVueStore,
).Build()

var UseVueStore command.InlineHandler = func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*VueUseStoreConfig)
	name := strcase.ToCamel(config.Name)
	storeKey := strcase.ToKebab(config.Name)
	filename := fmt.Sprintf("use%sStore", name)

	b := new(bytes.Buffer)
	if err := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/vue/vue_use_store.gotmpl"), b)).
		AddTask(task.NewTemplateRenderTask("vue_use_store", map[string]string{"storeKey": storeKey, "name": name}, b, b)).
		AddTask(task.NewWriteFileToDiskTask(filename, ".ts", config.OutDir, false, b)).
		Run(); err != nil {
		return err
	}
	return nil
}
