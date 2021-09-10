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

type DataProvider = func() ([]byte, error)

type Transform = func([]byte) ([]byte, error)

type FileDesc struct {
	Path       string
	Data       DataProvider
	Overwrite  bool
	Transforms []Transform
}

func (fd *FileDesc) String() string {
	return fmt.Sprintf("[FileDesc] path=%s overwrite=%v transforms=%d", fd.Path, fd.Overwrite, len(fd.Transforms))
}

type FileGenerator struct {
	file       string // file path to be generated
	data       DataProvider
	overwrite  bool
	transforms []Transform
}

func NewFileGenerator(fd *FileDesc) Runner {
	file := strings.TrimPrefix(fd.Path, "/")
	return &FileGenerator{file, fd.Data, fd.Overwrite, fd.Transforms}
}

func (fg *FileGenerator) Run() error {
	file := fg.file
	if file == "" {
		log.Warn("no path specified [skipped]")
		return nil
	}

	provider := fg.data
	if provider == nil {
		log.WithField("file", file).Warnf("no provider found [skipped]")
		return nil
	}

	overwrite := fg.overwrite
	if !overwrite {
		if _, err := os.Stat(file); err == nil {
			log.WithFields(log.Fields{"file": file, "overwrite": overwrite}).Info("emit file [skipped]")
			return nil
		}
	}

	data, err := provider()
	if err != nil {
		return errors.Wrapf(err, "failed to get from provider [ %s ]", file)
	}

	for _, t := range fg.transforms {
		data, err = t(data)
		if err != nil {
			return errors.Wrapf(err, "failed to apply transformation to [ %s ]", file)
		}
	}

	err = fg.write(file, data)
	if err == nil {
		log.WithFields(log.Fields{
			"overwrite": overwrite,
			"file":      fg.file,
		}).Infof("emit file")
	}
	return errors.Wrapf(err, "failed to generate [ %s ]", file)
}

func (fg *FileGenerator) write(file string, data []byte) error {
	if dir := path.Dir(file); dir != "" {
		if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to make '%s'", dir)
		}
	}
	return errors.Wrapf(os.WriteFile(file, data, fs.ModePerm), "failed to write file '%s'", file)
}
