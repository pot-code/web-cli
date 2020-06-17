package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Clone Git clone command
func Clone(ctx context.Context, url string, opts []string) ([]byte, error) {
	// check initial done
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	args := append([]string{"clone", url}, opts...)
	c := exec.CommandContext(ctx, "git", args...)
	console, err := c.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return console, fmt.Errorf("Git is not installed(https://git-scm.com/) or not exists in PATH")
		}
		return console, fmt.Errorf("Execution failed: %w", err)
	}
	return console, nil
}

// Goreturns run goreturn in given directory
func Goreturns(ctx context.Context, cwd string) ([]byte, error) {
	// check initial done
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	args := []string{"-w", "."}
	c := exec.CommandContext(ctx, "goreturns", args...)
	c.Dir = cwd
	console, err := c.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return console, fmt.Errorf("goreturns is not installed(go install github.com/sqs/goreturns) or not exists in PATH")
		}
		return console, fmt.Errorf("Execution failed: %w", err)
	}
	return console, nil
}

// Goimports run goimports
func Goimports(ctx context.Context, cwd string) ([]byte, error) {
	// check initial done
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	args := []string{"-w", "."}
	c := exec.CommandContext(ctx, "goimports", args...)
	c.Dir = cwd
	console, err := c.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return console, fmt.Errorf("goimports is not installed(go install github.com/bradfitz/goimports) or not exists in PATH")
		}
		return console, fmt.Errorf("Execution failed: %w", err)
	}
	return console, nil
}

// GoModTidy run go mod tidy
func GoModTidy(ctx context.Context, cwd string) ([]byte, error) {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	c := exec.CommandContext(ctx, "go", "mod", "tidy")
	c.Dir = cwd
	console, err := c.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return console, fmt.Errorf("Go is not installed(https://golang.org/dl/) or not exists in PATH")
		}
		return console, fmt.Errorf("Execution failed: %w", err)
	}
	return console, nil
}

// GoModInit init go module
func GoModInit(ctx context.Context, moduleName, cwd string) ([]byte, error) {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	c := exec.CommandContext(ctx, "go", "mod", "init", moduleName)
	c.Dir = cwd
	console, err := c.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return console, fmt.Errorf("Go is not installed(https://golang.org/dl/) or not exists in PATH")
		}
		return console, fmt.Errorf("Execution failed: %w", err)
	}
	return console, nil
}

// GoToCobraMap map from go type to cobra
var GoToCobraMap = map[string]string{
	"string":   "StringVar",
	"int":      "IntVar",
	"float":    "Float64Var",
	"duration": "DurationVar",
	"[]int":    "IntSliceVar",
	"[]string": "StringSliceVar",
}

// GoZeroValueMap string representation of zero value
var GoZeroValueMap = map[string]string{
	"string":   `""`,
	"int":      "0",
	"float":    "0.0",
	"duration": "0",
	"[]int":    "nil",
	"[]string": "nil",
}

// GoTypeToCobra returns cobra type with respect to golang type t
func GoTypeToCobra(t string) (string, error) {
	if to, ok := GoToCobraMap[t]; ok {
		return to, nil
	}
	return "", fmt.Errorf("Unsupported type '%s'", t)
}

// GetValueString returns string representation of golang value
func GetValueString(t string, val interface{}) (string, error) {
	if val == nil {
		if zero, ok := GoZeroValueMap[t]; ok {
			return zero, nil
		}
		return "", fmt.Errorf("Unsupported type '%s'", t)
	}
	switch t {
	case "string":
		return fmt.Sprintf(`"%s"`, val.(string)), nil
	case "int":
		return fmt.Sprintf("%d", val.(int)), nil
	case "duration":
		return fmt.Sprintf("%d", val.(time.Duration)), nil
	case "float":
		return fmt.Sprintf("%f", val.(float64)), nil
	case "[]int":
		strs := make([]string, len(val.([]int)))
		for i, v := range val.([]int) {
			strs[i] = strconv.Itoa(v)
		}
		return fmt.Sprintf("[]int{%s}", strings.Join(strs, ",")), nil
	case "[]string":
		return fmt.Sprintf(`[]string{"%s"}`, strings.Join(val.([]string), `","`)), nil
	}
	return "", fmt.Errorf("Unsuppoted type '%s'", t)
}

// ToKebabCase transform var name to kebab case
//
// eg:
// IDString -> id-string
// HelloWorld -> hello-world
// appName -> app-name
// prefixName -> prefix-name
func ToKebabCase(s string) (string, error) {
	if s == "" {
		return s, nil
	}
	validate := regexp.MustCompile("^[A-Za-z]+$")
	if !validate.MatchString(s) {
		return s, fmt.Errorf("Malformed input string '%s'", s)
	}

	splitPosition := []int{0}
	lastUpper := unicode.IsUpper(rune(s[0]))
	for i := 1; i < len(s); i++ {
		r := rune(s[i])
		if unicode.IsLower(r) {
			if lastUpper && i > 1 {
				splitPosition = append(splitPosition, i-1)
			}
			lastUpper = false
		} else {
			lastUpper = true
		}
	}
	if len(splitPosition) < 2 {
		return strings.ToLower(s), nil
	}
	// add a length value for final segment
	splitPosition = append(splitPosition, len(s))
	var parts []string
	for i := 1; i < len(splitPosition); i++ {
		parts = append(parts, strings.ToLower(s[splitPosition[i-1]:splitPosition[i]]))
	}
	return strings.Join(parts, "-"), nil
}
