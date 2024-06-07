package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
)

var variableNameExp = `^[-_\w][-_\w\d]+$`
var variableNameReg = regexp.MustCompile(variableNameExp)

func ValidateVariableName(name string) error {
	if !variableNameReg.MatchString(name) {
		return fmt.Errorf("input must be in form: %s", variableNameExp)
	}
	return nil
}

func ValidateVersion(v string) error {
	if _, err := version.NewVersion(v); err != nil {
		return fmt.Errorf("invalid version format")
	}
	return nil
}
