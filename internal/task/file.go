package task

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type WriteFileToDiskTask struct {
	Name      string
	Suffix    string
	Folder    string
	Overwrite bool
	data      io.Reader
}

func NewWriteFileToDiskTask(name string, suffix string, folder string, overwrite bool, data io.Reader) *WriteFileToDiskTask {
	return &WriteFileToDiskTask{
		Name:      name,
		Suffix:    suffix,
		Folder:    folder,
		Overwrite: overwrite,
		data:      data,
	}
}

func (wft *WriteFileToDiskTask) Run() error {
	if wft.shouldSkip() {
		log.WithFields(log.Fields{
			"file":      wft.getFullPath(),
			"overwrite": wft.Overwrite,
		}).Info("emit file [skipped]")
		return nil
	}

	data, err := ioutil.ReadAll(wft.data)
	if err != nil {
		return fmt.Errorf("WriteFileTask: read data from data source: %w", err)
	}

	if err := wft.write(data); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"file": wft.getFullPath(),
	}).Info("emit file")
	return nil
}

func (wft *WriteFileToDiskTask) shouldSkip() bool {
	return fileExists(wft.getFullPath()) && !wft.Overwrite
}

func (wft *WriteFileToDiskTask) getFullPath() string {
	return path.Join(wft.Folder, wft.Name+wft.Suffix)
}

func (wft *WriteFileToDiskTask) write(data []byte) error {
	if err := wft.mkdirIfNecessary(); err != nil {
		return err
	}

	filePath := wft.getFullPath()
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("open file [path: %s]: %w", filePath, err)
	}
	defer fd.Close()

	n, err := fd.Write(data)
	if err != nil {
		return errors.Wrapf(err, "failed to write data to '%s'", filePath)
	}
	log.WithFields(log.Fields{"bytes": n, "file": filePath}).Debug("write file")
	return nil
}

func (wft *WriteFileToDiskTask) mkdirIfNecessary() error {
	dir := wft.Folder
	if dir == "" {
		return nil
	}
	if fileExists(dir) {
		return nil
	}
	return errors.Wrapf(os.MkdirAll(dir, fs.ModePerm), "failed to make '%s'", dir)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
