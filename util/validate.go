package util

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
)

var NameRegexp = `^[_\w][-_\w]*$`

func ValidateProjectName(name string) error {
	reg := regexp.MustCompile(NameRegexp)
	if !reg.MatchString(name) {
		return fmt.Errorf("input must be adhere to %s", NameRegexp)
	}
	return nil
}

func ValidateUserName(name string) error {
	reg := regexp.MustCompile(NameRegexp)
	if !reg.MatchString(name) {
		return fmt.Errorf("input must be adhere to %s", NameRegexp)
	}
	return nil
}

func ValidateVersion(v string) error {
	if _, err := version.NewVersion(v); err != nil {
		return fmt.Errorf("invalid version format")
	}
	return nil
}
