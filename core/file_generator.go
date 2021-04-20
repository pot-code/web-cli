package core

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DataProvider = func() []byte

type FileGenerator struct {
	file    string // file path to be generated
	data    DataProvider
	cleaned bool
}

func NewFileGenerator(path string, data DataProvider) Generator {
	file := strings.TrimPrefix(path, "/")
	log.Trace("registered file: ", file)
	return &FileGenerator{file, data, false}
}

func (gt *FileGenerator) Run() error {
	file := gt.file
	provider := gt.data

	if file == "" {
		log.Info("[skipped]no path specified")
		return nil
	}

	if provider == nil {
		log.Infof("[skipped]no provider for '%s'", gt.file)
		return nil
	}
	err := gt.write(file, provider())
	if err == nil {
		log.Infof("emit '%s'", gt.file)
	}
	return errors.Wrapf(err, "failed to generate '%s'", file)
}

func (gt *FileGenerator) Cleanup() error {
	if gt.cleaned {
		return nil
	}
	gt.cleaned = true

	log.Debugf("removing file '%s'", gt.file)
	err := os.Remove(gt.file)
	if err != nil {
		if !os.IsNotExist(err) {
			log.WithFields(log.Fields{"error": err.Error(), "file": gt.file}).Debug("[cleanup]failed to cleanup")
		}
	}

	return errors.WithStack(err)
}

func (gt *FileGenerator) write(file string, data []byte) error {
	if dir := path.Dir(file); dir != "" {
		if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to make '%s'", dir)
		}
	}
	return errors.Wrapf(os.WriteFile(file, data, fs.ModePerm), "failed to write file '%s'", file)
}

func (gt *FileGenerator) String() string {
	return fmt.Sprintf("[FileGenerator]path=%s", gt.file)
}
