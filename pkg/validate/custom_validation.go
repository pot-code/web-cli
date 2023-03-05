package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
)

var naturalNameExp = `^[-_\w][-_\w\d]+$`
var natureNameReg = regexp.MustCompile(naturalNameExp)

func ValidateNatureName(name string) error {
	if !natureNameReg.MatchString(name) {
		return fmt.Errorf("input must be in form: %s", naturalNameExp)
	}
	return nil
}

func ValidateVersion(v string) error {
	if _, err := version.NewVersion(v); err != nil {
		return fmt.Errorf("invalid version format")
	}
	return nil
}
