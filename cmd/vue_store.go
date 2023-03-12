package cmd

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/provider"
	"github.com/pot-code/web-cli/pkg/task"
	"github.com/urfave/cli/v2"
)

const (
	VueStoreSuffix              = ".ts"
	VueUseStoreFileNameTemplate = "use%sStore"
)

type VueUseStoreConfig struct {
	OutDir string `flag:"output" alias:"o" usage:"destination directory"`
	Name   string `arg:"0" alias:"MODULE_NAME" validate:"required,var"`
}

var VueUseStoreCmd = command.NewCliCommand("store", "add vue pinia store",
	new(VueUseStoreConfig),
	command.WithArgUsage("MODULE_NAME"),
	command.WithAlias([]string{"s"}),
).AddHandlers(
	UseVueStore,
).BuildCommand()

var UseVueStore command.InlineHandler = func(c *cli.Context, cfg interface{}) error {
	config := cfg.(*VueUseStoreConfig)
	moduleName := strcase.ToCamel(config.Name)
	storeName := strcase.ToKebab(config.Name)
	fileName := fmt.Sprintf(VueUseStoreFileNameTemplate, moduleName)

	b := new(bytes.Buffer)
	if err := task.NewSequentialScheduler().
		AddTask(task.NewReadFromProviderTask(provider.NewEmbedFileProvider("templates/vue_use_store.gohtml"), b)).
		AddTask(task.NewTemplateRenderTask("vue_use_store", map[string]string{"moduleName": moduleName, "storeName": storeName}, b, b)).
		AddTask(task.NewWriteFileToDiskTask(fileName, VueStoreSuffix, config.OutDir, false, b)).
		Run(); err != nil {
		return err
	}
	return nil
}
