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

type FileDesc struct {
	Path string
	Data DataProvider
}

func (fd FileDesc) String() string {
	return fd.Path
}

type DataProvider = func() []byte

type FileGenerator struct {
	file    string // file path to be generated
	data    DataProvider
	cleaned bool
}

func NewFileGenerator(fd *FileDesc) Generator {
	file := strings.TrimPrefix(fd.Path, "/")
	log.Trace("registered file: ", file)
	return &FileGenerator{file, fd.Data, false}
}

func (fg *FileGenerator) Run() error {
	file := fg.file
	provider := fg.data

	if file == "" {
		log.Info("[skipped]no path specified")
		return nil
	}

	if provider == nil {
		log.Infof("[skipped]no provider for '%s'", fg.file)
		return nil
	}
	err := fg.write(file, provider())
	if err == nil {
		log.Infof("emit '%s'", fg.file)
	}
	return errors.Wrapf(err, "failed to generate '%s'", file)
}

func (fg *FileGenerator) Cleanup() error {
	if fg.cleaned {
		return nil
	}
	fg.cleaned = true

	log.Debugf("removing file '%s'", fg.file)
	err := os.Remove(fg.file)
	if err != nil {
		if !os.IsNotExist(err) {
			log.WithFields(log.Fields{"error": err.Error(), "file": fg.file}).Debug("[cleanup]failed to cleanup")
		}
	}

	return errors.WithStack(err)
}

func (fg *FileGenerator) write(file string, data []byte) error {
	if dir := path.Dir(file); dir != "" {
		if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to make '%s'", dir)
		}
	}
	return errors.Wrapf(os.WriteFile(file, data, fs.ModePerm), "failed to write file '%s'", file)
}

func (fg *FileGenerator) String() string {
	return fmt.Sprintf("[FileGenerator]path=%s", fg.file)
}
