package task

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type SequentialExecutor struct {
	tasks []Task
}

var _ Task = &SequentialExecutor{}

func NewSequentialExecutor(tasks ...Task) *SequentialExecutor {
	return &SequentialExecutor{tasks}
}

func (se SequentialExecutor) Run() error {
	start := time.Now()

	for i, t := range se.tasks {
		log.WithFields(log.Fields{
			"number": i,
			"total":  len(se.tasks),
		}).Debug("SequentialExecutor running task")

		if err := t.Run(); err != nil {
			return errors.Wrap(err, "failed to run task")
		}
	}

	log.WithFields(log.Fields{
		"duration": time.Since(start),
		"total":    len(se.tasks),
	}).Info("SequentialExecutor finished")
	return nil
}
