package task

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

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

	if err := wft.write(); err != nil {
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

func (wft *WriteFileToDiskTask) write() error {
	if err := wft.mkdirIfNecessary(); err != nil {
		return err
	}

	filePath := wft.getFullPath()
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("open file [path: %s]: %w", filePath, err)
	}
	defer fd.Close()

	n, err := io.Copy(fd, wft.data)
	if err != nil {
		return fmt.Errorf("write data to %s: %w", filePath, err)
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
	err := os.MkdirAll(dir, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("make dir %s: %w", dir, err)
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type FileProvider interface {
	Get() (io.Reader, error)
}

type GenerateFileFromProviderTask struct {
	fileName  string
	suffix    string
	folder    string
	overwrite bool
	provider  FileProvider
}

func NewGenerateFileFromProvider(
	fileName string,
	suffix string,
	folder string,
	overwrite bool,
	provider FileProvider,
) *GenerateFileFromProviderTask {
	return &GenerateFileFromProviderTask{fileName: fileName, suffix: suffix, folder: folder, overwrite: overwrite, provider: provider}
}

func (t *GenerateFileFromProviderTask) Run() error {
	fd, err := t.provider.Get()
	if err != nil {
		return fmt.Errorf("get template from provider: %w", err)
	}

	fdt := NewWriteFileToDiskTask(t.fileName, t.suffix, t.folder, t.overwrite, fd)
	if err := fdt.Run(); err != nil {
		return fmt.Errorf("write file to disk [file_name: %s folder: %s]: %w", t.fileName, t.folder, err)
	}
	return nil
}
