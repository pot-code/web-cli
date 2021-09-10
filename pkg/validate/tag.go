package validate

import (
	"reflect"
	"strings"
)

// GetOneOfItems parse and get oneof items from field tag
func GetOneOfItems(tag reflect.StructTag) []string {
	if tag == "" {
		return nil
	}

	vt := tag.Get("validate")
	if vt == "" {
		return nil
	}

	options := strings.Split(vt, ",")
	for _, o := range options {
		if strings.HasPrefix(o, "oneof") {
			return strings.Split(strings.TrimPrefix(o, "oneof="), " ")
		}
	}
	return nil
}

// IsRequired check if tag has required flag in validate tag
func IsRequired(tag reflect.StructTag) bool {
	if tag == "" {
		return false
	}

	vt := tag.Get("validate")
	if vt == "" {
		return false
	}

	return strings.Contains(vt, "required")
}
