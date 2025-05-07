package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
)

var identifierRegExp = `^[-_\w][-_\w\d]+$`
var identifierReg = regexp.MustCompile(identifierRegExp)

// validateIdentifier 校验标识符类型字符串
func validateIdentifier(name string) error {
	if !identifierReg.MatchString(name) {
		return fmt.Errorf("input must be in form: %s", identifierRegExp)
	}
	return nil
}

// validateVersion 校验版本号
func validateVersion(v string) error {
	if _, err := version.NewVersion(v); err != nil {
		return fmt.Errorf("invalid version format")
	}
	return nil
}
