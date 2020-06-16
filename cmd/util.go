package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
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
	// stderr, err := c.StderrPipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to bind stderr pipe: %w", err)
	// }

	// err = c.Start()
	// if err != nil {
	// 	if errors.Is(err, exec.ErrNotFound) {
	// 		return nil, fmt.Errorf("Git is not installed(https://git-scm.com/) or not exists in PATH")
	// 	}
	// 	return nil, fmt.Errorf("Failed to start command: %w", err)
	// }
	// console, err := readConsoleIntoBuffer(stderr)
	// if err != nil {
	// 	return nil, err
	// }

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
