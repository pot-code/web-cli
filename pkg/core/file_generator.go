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
	Path      string
	Data      DataProvider
	Overwrite bool
}

func (fd FileDesc) String() string {
	return fd.Path
}

type DataProvider = func() []byte

type FileGenerator struct {
	file      string // file path to be generated
	data      DataProvider
	cleaned   bool
	overwrite bool
}

func NewFileGenerator(fd *FileDesc) Runner {
	file := strings.TrimPrefix(fd.Path, "/")
	log.Trace("registered file: ", file)
	return &FileGenerator{file, fd.Data, false, fd.Overwrite}
}

func (fg *FileGenerator) Run() error {
	file := fg.file
	provider := fg.data
	overwrite := fg.overwrite

	if file == "" {
		log.Info("[skip]no path specified")
		return nil
	}

	if provider == nil {
		log.Infof("[skip]no provider for '%s'", fg.file)
		return nil
	}

	if !overwrite {
		if _, err := os.Stat(file); err == nil {
			log.WithFields(log.Fields{"file": file, "overwrite": overwrite}).Info("[skip]no overwrite")
			return nil
		}
	}

	err := fg.write(file, provider())
	if err == nil {
		log.WithFields(log.Fields{
			"overwrite": overwrite,
		}).Infof("emit '%s'", fg.file)
	}
	return errors.Wrapf(err, "failed to generate '%s'", file)
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
