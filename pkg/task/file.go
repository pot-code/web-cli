package task

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DataSource = func(buf *bytes.Buffer) error

type Pipeline = func(src *bytes.Buffer, dst *bytes.Buffer) error

type FileDesc struct {
	Path       string
	Source     DataSource
	Overwrite  bool
	Transforms []Pipeline
}

func (fd *FileDesc) String() string {
	return fmt.Sprintf("path=%s overwrite=%v transforms=%d", fd.Path, fd.Overwrite, len(fd.Transforms))
}

type FileGenerator struct {
	file string // file path to be generated
	fd   *FileDesc
}

func NewFileGenerator(fd *FileDesc) Runner {
	file := strings.TrimPrefix(fd.Path, "/")
	return &FileGenerator{file, fd}
}

func (fg *FileGenerator) Run() error {
	file := fg.file
	if file == "" {
		log.Warn("no path specified [skipped]")
		return nil
	}

	Source := fg.fd.Source
	if Source == nil {
		log.WithField("file", file).Warnf("no provider found [skipped]")
		return nil
	}

	overwrite := fg.fd.Overwrite
	if !overwrite {
		if _, err := os.Stat(file); err == nil {
			log.WithFields(log.Fields{"file": file, "overwrite": overwrite}).Info("emit file [skipped]")
			return nil
		}
	}

	buf := new(bytes.Buffer)
	err := Source(buf)
	if err != nil {
		return errors.Wrapf(err, "failed to get data from provider [ %s ]", file)
	}

	for _, t := range fg.fd.Transforms {
		out := new(bytes.Buffer)
		err = t(buf, out)
		if err != nil {
			return errors.Wrapf(err, "failed to apply transformation to [ %s ]", file)
		}
		buf = out
	}

	err = fg.write(file, buf.Bytes())
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
