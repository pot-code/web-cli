package task

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Task interface {
	Run() error
}

type SequentialExecutor struct {
	tasks []Task
}

var _ Task = (*SequentialExecutor)(nil)

func NewSequentialExecutor() *SequentialExecutor {
	return &SequentialExecutor{}
}

func (se *SequentialExecutor) AddTask(task Task) *SequentialExecutor {
	se.tasks = append(se.tasks, task)
	return se
}

func (se *SequentialExecutor) Run() error {
	start := time.Now()
	total := len(se.tasks)

	log.WithFields(log.Fields{"task_total": total}).Debug("SequentialExecutor start")
	for i, t := range se.tasks {
		log.WithFields(log.Fields{"task_total": total}).Debugf("SequentialExecutor running task #%d", i+1)

		if err := t.Run(); err != nil {
			return fmt.Errorf("run task: %w", err)
		}
	}

	log.WithFields(log.Fields{"duration": time.Since(start)}).Debug("SequentialExecutor finished")
	return nil
}

type ParallelExecutor struct {
	tasks []Task
}

var _ Task = (*ParallelExecutor)(nil)

func NewParallelExecutor() *ParallelExecutor {
	return &ParallelExecutor{}
}

func (pe *ParallelExecutor) AddTask(task Task) *ParallelExecutor {
	pe.tasks = append(pe.tasks, task)
	return pe
}

func (pe *ParallelExecutor) Run() error {
	start := time.Now()
	total := len(pe.tasks)
	eg := new(errgroup.Group)

	log.WithFields(log.Fields{"task_total": total}).Debug("ParallelExecutor start")
	for _, t := range pe.tasks {
		task := t // fix loopclosure
		eg.Go(func() error {
			return task.Run()
		})
	}
	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("run task: %w", err)
	}

	log.WithFields(log.Fields{
		"duration": time.Since(start),
		"total":    len(pe.tasks),
	}).Debug("ParallelExecutor finished")
	return nil
}
