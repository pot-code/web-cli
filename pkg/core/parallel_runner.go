package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ParallelRunner struct {
	commands []Executor
	files    []Generator
	cleaned  bool
}

var _ Generator = ParallelRunner{}

func NewParallelRunner(tasks ...Executor) *ParallelRunner {
	var (
		commands []Executor
		files    []Generator
	)

	for _, t := range tasks {
		if gen, ok := t.(Generator); ok {
			files = append(files, gen)
		} else {
			commands = append(commands, t)
		}
	}
	return &ParallelRunner{commands, files, false}
}

func (pg ParallelRunner) Run() error {
	start := time.Now()

	if err := pg.runCommands(); err != nil {
		return err
	}

	if err := pg.runGenerators(); err != nil {
		return err
	}

	log.WithField("duration", time.Since(start)).Info("parallel runner finished")
	return nil
}

func (pg ParallelRunner) runCommands() error {
	for _, cmd := range pg.commands {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (pg ParallelRunner) runGenerators() error {
	files := pg.files
	errChan := make(chan error)
	doneChan := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := range files {
		wg.Add(1)
		go func(task Generator) {
			if err := task.Run(); err != nil {
				errChan <- task.Run()
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

func (pg ParallelRunner) Cleanup() error {
	if pg.cleaned {
		return nil
	}
	pg.cleaned = true

	for _, task := range pg.files {
		if ce := task.Cleanup(); ce != nil {
			return errors.WithStack(ce)
		}
	}

	return nil
}

func (pg ParallelRunner) String() string {
	return fmt.Sprintf("[ParallelRunner]tasks=%d", len(pg.commands))
}
