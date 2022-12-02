package env

import (
	"os"
	"path"
)

const TemplatePathEnv = "WEB_CLI_TEMPLATE_PATH"

func GetAbsoluteTemplatePath(name string) string {
	e, ok := os.LookupEnv("WEB_CLI_TEMPLATE_PATH")
	if ok {
		return path.Join(e, name)
	}
	return ""
}
