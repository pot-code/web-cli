package task

import (
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Task interface {
	Run() error
}

type SequentialScheduler struct {
	tasks []Task
}

var _ Task = (*SequentialScheduler)(nil)

func NewSequentialScheduler() *SequentialScheduler {
	return &SequentialScheduler{}
}

func (se *SequentialScheduler) AddTask(task Task) *SequentialScheduler {
	se.tasks = append(se.tasks, task)
	return se
}

func (se *SequentialScheduler) Run() error {
	start := time.Now()
	total := len(se.tasks)

	log.WithFields(log.Fields{"task_total": total}).Debug("SequentialScheduler start")
	for i, t := range se.tasks {
		log.WithFields(log.Fields{"task_total": total}).Debugf("SequentialScheduler running task #%d", i+1)

		if err := t.Run(); err != nil {
			return err
		}
	}

	log.WithFields(log.Fields{"duration": time.Since(start)}).Debug("SequentialScheduler finished")
	return nil
}

type ParallelScheduler struct {
	tasks []Task
}

var _ Task = (*ParallelScheduler)(nil)

func NewParallelScheduler() *ParallelScheduler {
	return &ParallelScheduler{}
}

func (pe *ParallelScheduler) AddTask(task Task) *ParallelScheduler {
	pe.tasks = append(pe.tasks, task)
	return pe
}

func (pe *ParallelScheduler) Run() error {
	start := time.Now()
	total := len(pe.tasks)
	eg := new(errgroup.Group)

	log.WithFields(log.Fields{"task_total": total}).Debug("ParallelScheduler start")
	for _, t := range pe.tasks {
		task := t // fix loop closure
		eg.Go(func() error {
			return task.Run()
		})
	}
	err := eg.Wait()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"duration": time.Since(start),
		"total":    len(pe.tasks),
	}).Debug("ParallelScheduler finished")
	return nil
}
