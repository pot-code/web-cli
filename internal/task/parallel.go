package task

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type ParallelExecutor struct {
	tasks []Task
}

var _ Task = &ParallelExecutor{}

func NewParallelExecutor(tasks ...Task) *ParallelExecutor {
	return &ParallelExecutor{tasks}
}

func (pe ParallelExecutor) Run() error {
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
	log.WithFields(log.Fields{
		"duration": time.Since(start),
		"total":    len(pe.tasks),
	}).Debug("ParallelExecutor finished")
	return errors.Wrap(err, "failed to run task")
}
