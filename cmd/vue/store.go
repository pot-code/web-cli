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

var VueUseStoreCmd = command.NewCommand("store", "add vue pinia store",
	new(VueUseStoreConfig),
	command.WithAlias([]string{"s"}),
).AddHandler(
	UseVueStore,
).Create()

var UseVueStore = command.InlineHandler[*VueUseStoreConfig](func(c *cli.Context, config *VueUseStoreConfig) error {
	name := strcase.ToCamel(config.Name)
	storeKey := strcase.ToKebab(config.Name)
	filename := fmt.Sprintf("use%sStore", name)

	b := new(bytes.Buffer)
	if err := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/vue/vue_use_store.go.tmpl"), b)).
		AddTask(task.NewTemplateRenderTask("vue_use_store", map[string]string{"storeKey": storeKey, "name": name}, b, b)).
		AddTask(task.NewWriteFileToDiskTask(filename, ".ts", config.OutDir, false, b)).
		Run(); err != nil {
		return err
	}
	return nil
})
