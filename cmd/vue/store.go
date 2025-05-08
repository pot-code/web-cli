package vue

import (
	"bytes"
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

type VueUseStoreConfig struct {
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
	Name   string `arg:"0" alias:"MODULE_NAME" validate:"required,identifier"`
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
		AddTask(task.NewWriteFileToDiskTask(filename, ".ts", b, task.WithFolder(config.OutDir))).
		Run(); err != nil {
		return err
	}
	return nil
})
