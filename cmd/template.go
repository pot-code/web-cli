package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"text/template"

	"gopkg.in/yaml.v2"
)

// TemplateConfigItem entry in AppConfig struct
type TemplateConfigItem struct {
	Name         string `yaml:"name"`     // key name
	DefaultValue string `yaml:"value"`    // default field value
	Usage        string `yaml:"usage"`    // usage string, can be treat as the app name
	Required     bool   `yaml:"required"` // is field required
}

// TemplateEntry TODO
type TemplateEntry struct {
	Template string // template string
	Name     string // template name
	Output   string // output path
}

// TemplateData data struct used for generating template
type TemplateData struct {
	EnvPrefix      string                `yaml:"envPrefix"`  // prefix of env variables for lookup
	ModuleName     string                `yaml:"moduleName"` // go module name used in go.mod
	Usage          string                `yaml:"usage"`      // application usage string
	Short          string                `yaml:"short"`      // application short desc string
	ConfigGlobal   []*TemplateConfigItem `yaml:"global"`
	ConfigLogging  []*TemplateConfigItem `yaml:"logging"`
	ConfigSecurity []*TemplateConfigItem `yaml:"security"`
}

// NewTemplateData create new TemplateModel
func NewTemplateData(envPrefix, moduleName string, usage string) *TemplateData {
	return &TemplateData{
		EnvPrefix:  envPrefix,
		ModuleName: moduleName,
		Usage:      usage,
	}
}

// LoadFromYaml read field data from yaml file
func (tm *TemplateData) LoadFromYaml(path string) error {
	fd, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return fmt.Errorf("Can't read yaml config: %w", err)
	}
	decoder := yaml.NewDecoder(fd)
	if err = decoder.Decode(tm); err != nil {
		return fmt.Errorf("Failed to decode yaml: %w", err)
	}
	return nil
}

// CreateFromTemplate creates file from template
//
// content is template string, name is template name, out is output path
func CreateFromTemplate(entry *TemplateEntry, data interface{}) (err error) {
	out := entry.Output
	dir := path.Dir(out)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("Error while creating directories: %w", err)
	}
	fd, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0666)
	defer fd.Close()
	if err != nil {
		return err
	}
	// delete the output file if any error occurs
	defer func() {
		if err != nil {
			if err = os.Remove(out); err != nil {
				log.Printf("Failed to cleanup: %s", err)
			}
		}
	}()
	parsed, err := template.New(entry.Name).Parse(entry.Template)
	if err != nil {
		return fmt.Errorf("Error while parse the template '%s': %w", entry.Name, err)
	}
	return parsed.Execute(fd, data)
}
