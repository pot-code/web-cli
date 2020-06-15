package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

// Clone Git clone command
func Clone(ctx context.Context, url string, opts []string) error {
	// check initial done
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	args := append([]string{"clone", url}, opts...)
	c := exec.CommandContext(ctx, "git", args...)
	stderr, err := c.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to bind stdout pipe: %w", err)
	}
	outScanner := bufio.NewScanner(stderr)
	err = c.Start()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("Git is not installed(https://git-scm.com/) or not exists in PATH")
		}
		return fmt.Errorf("Failed to start command: %w", err)
	}
	for outScanner.Scan() {
		log.Printf("[%s][%s]: %s", c.Path, url, outScanner.Text())
	}

	// if ctx != nil {
	// 	waitDone := make(chan struct{})
	// 	go func() {
	// 		err = c.Wait()
	// 		close(waitDone)
	// 	}()
	// 	select {
	// 	case <-ctx.Done():
	// 		c.Process.Kill()
	// 		return ctx.Err()
	// 	case <-waitDone:
	// 	}
	// } else {
	// 	err = c.Wait()
	// }
	if err = c.Wait(); err != nil {
		return fmt.Errorf("Execution failed: %w", err)
	}
	return nil
}

// Goreturns run goreturn in given directory
func Goreturns(ctx context.Context, dir string) error {
	// check initial done
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	args := []string{"-w", "."}
	c := exec.CommandContext(ctx, "goreturns", args...)
	c.Dir = dir
	stderr, err := c.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to bind stdout pipe: %w", err)
	}
	outScanner := bufio.NewScanner(stderr)
	err = c.Start()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("goreturns is not installed(go install github.com/sqs/goreturns) or not exists in PATH")
		}
		return fmt.Errorf("Failed to start command: %w", err)
	}
	for outScanner.Scan() {
		log.Printf("[%s]: %s", c.Path, outScanner.Text())
	}

	if err = c.Wait(); err != nil {
		return fmt.Errorf("Execution failed: %w", err)
	}
	return nil
}

// Goimports run goimports
func Goimports(ctx context.Context, dir string) error {
	// check initial done
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	args := []string{"-w", "."}
	c := exec.CommandContext(ctx, "goimports", args...)
	c.Dir = dir
	stderr, err := c.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to bind stdout pipe: %w", err)
	}
	outScanner := bufio.NewScanner(stderr)
	err = c.Start()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("goimports is not installed(go install github.com/bradfitz/goimports) or not exists in PATH")
		}
		return fmt.Errorf("Failed to start command: %w", err)
	}
	for outScanner.Scan() {
		log.Printf("[%s]: %s", c.Path, outScanner.Text())
	}

	if err = c.Wait(); err != nil {
		return fmt.Errorf("Execution failed: %w", err)
	}
	return nil
}
