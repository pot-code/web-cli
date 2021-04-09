package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ParallelGenerator struct {
	subtasks []Generator
	cleaned  bool
}

var _ Generator = ParallelGenerator{}

func NewParallelGenerator(subtasks ...Generator) *ParallelGenerator {
	return &ParallelGenerator{subtasks, false}
}

func (pg ParallelGenerator) Gen() error {
	subtasks := pg.subtasks
	errChan := make(chan error)
	doneChan := make(chan struct{})
	start := time.Now()
	wg := sync.WaitGroup{}

	for i := range subtasks {
		wg.Add(1)
		go func(task Generator) {
			if err := task.Gen(); err != nil {
				errChan <- task.Gen()
			} else {
				wg.Done()
			}
		}(subtasks[i])
	}
	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		log.WithField("duration", time.Since(start)).Info("module generation finished")
		return nil
	}
}

func (pg ParallelGenerator) Cleanup() error {
	if pg.cleaned {
		return nil
	}
	pg.cleaned = true

	for _, task := range pg.subtasks {
		if ce := task.Cleanup(); ce != nil {
			return errors.WithStack(ce)
		}
	}

	return nil
}

func (pg ParallelGenerator) String() string {
	return fmt.Sprintf("[ModuleGenerator]tasks=%d", len(pg.subtasks))
}
