package core

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
	eg := new(errgroup.Group)

	for _, t := range files {
		task := t // fix loopclosure
		eg.Go(func() error {
			return task.Run()
		})
	}
	return errors.Wrap(eg.Wait(), "failed to run generators")
}
