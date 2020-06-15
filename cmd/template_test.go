package cmd

import (
	"testing"

	"github.com/pot-code/web-cli/template/backend"
)

func TestConfigTemplate(t *testing.T) {
	model := NewTemplateData("TEST_APP", "github.com/pot-code/test", "app")

	if err := CreateFromTemplate(&TemplateEntry{
		Name:     "config",
		Output:   "./config.go",
		Template: backend.ConfigTemplate,
	}, model); err != nil {
		t.Fatal(err)
	}
}

func TestMainTemplate(t *testing.T) {
	model := NewTemplateData("TEST_APP", "github.com/pot-code/test", "app")

	if err := CreateFromTemplate(&TemplateEntry{
		Name:     "main",
		Output:   "./main.go",
		Template: backend.MainTemplate,
	}, model); err != nil {
		t.Fatal(err)
	}
}
