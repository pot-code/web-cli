package util

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
)

var NameExp = `^[_\w][-_\w]*$`

var NameReg = regexp.MustCompile(NameExp)

func ValidateVarName(name string) error {
	if !NameReg.MatchString(name) {
		return fmt.Errorf("input must be adhere to %s", NameExp)
	}
	return nil
}

func ValidateUserName(name string) error {
	if !NameReg.MatchString(name) {
		return fmt.Errorf("input must be adhere to %s", NameExp)
	}
	return nil
}

func ValidateVersion(v string) error {
	if _, err := version.NewVersion(v); err != nil {
		return fmt.Errorf("invalid version format")
	}
	return nil
}
