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
	total := len(se.tasks)

	log.WithFields(log.Fields{"task_total": total}).Debug("SequentialExecutor start")
	for i, t := range se.tasks {
		log.WithFields(log.Fields{"task_total": total}).Debugf("SequentialExecutor running task #%d", i+1)

		if err := t.Run(); err != nil {
			return errors.Wrap(err, "failed to run task")
		}
	}

	log.WithFields(log.Fields{"duration": time.Since(start)}).Debug("SequentialExecutor finished")
	return nil
}
