package core

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type ParallelRunner struct {
	beforeCommands []Runner
	commands       []Runner
	files          []Runner
	cleaned        bool
}

var _ Runner = &ParallelRunner{}

func NewParallelRunner(tasks ...Runner) *ParallelRunner {
	var (
		commands       []Runner
		beforeCommands []Runner
		files          []Runner
	)

	for _, t := range tasks {
		switch v := t.(type) {
		case *FileGenerator:
			files = append(files, v)
		case *ShellCmdExecutor:
			if v.cmd.Before {
				beforeCommands = append(beforeCommands, v)
			} else {
				commands = append(commands, v)
			}
		default:
			log.Fatalf("unknown task types")
		}
	}
	return &ParallelRunner{beforeCommands, commands, files, false}
}

func (pr ParallelRunner) Run() error {
	start := time.Now()

	if err := pr.runBeforeCommands(); err != nil {
		return err
	}

	if err := pr.runGenerators(); err != nil {
		return err
	}

	if err := pr.runCommands(); err != nil {
		return err
	}

	log.WithField("duration", time.Since(start)).Info("parallel runner finished")
	return nil
}

func (pr ParallelRunner) runCommands() error {
	for _, cmd := range pr.commands {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (pr ParallelRunner) runBeforeCommands() error {
	for _, cmd := range pr.beforeCommands {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (pr ParallelRunner) runGenerators() error {
	files := pr.files
	errChan := make(chan error)
	doneChan := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := range files {
		wg.Add(1)
		go func(task Runner) {
			if err := task.Run(); err != nil {
				errChan <- err
			} else {
				wg.Done()
			}
		}(files[i])
	}
	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

func (pr ParallelRunner) String() string {
	return fmt.Sprintf("[ParallelRunner]: generator=%d, executor=%d", len(pr.files), len(pr.commands))
}
