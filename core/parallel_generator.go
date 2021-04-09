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

func (mg ParallelGenerator) Gen() error {
	subtasks := mg.subtasks
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
		mg.Cleanup()
		return err
	case <-doneChan:
		log.WithField("duration", time.Since(start)).Info("module generation finished")
		return nil
	}
}

func (mg ParallelGenerator) Cleanup() error {
	if mg.cleaned {
		return nil
	}
	mg.cleaned = true

	for _, task := range mg.subtasks {
		if ce := task.Cleanup(); ce != nil {
			return errors.WithStack(ce)
		}
	}

	return nil
}

func (mg ParallelGenerator) String() string {
	return fmt.Sprintf("[ModuleGenerator]tasks=%d", len(mg.subtasks))
}
