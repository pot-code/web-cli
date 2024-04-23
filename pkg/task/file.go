package task

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

type WriteFileToDiskTask struct {
	// name filename
	name string
	// suffix file suffix
	suffix string
	folder string
	// overwrite if exists
	overwrite bool
	data      io.Reader
}

func NewWriteFileToDiskTask(name string, suffix string, folder string, overwrite bool, in io.Reader) *WriteFileToDiskTask {
	return &WriteFileToDiskTask{
		name:      name,
		suffix:    suffix,
		folder:    folder,
		overwrite: overwrite,
		data:      in,
	}
}

func (wft *WriteFileToDiskTask) Run() error {
	if wft.shouldSkip() {
		log.Info().Str("file", wft.getFullPath()).Bool("overwrite", wft.overwrite).Msg("emit file [skipped]")
		return nil
	}

	if err := wft.write(); err != nil {
		return err
	}

	log.Info().Str("file", wft.getFullPath()).Msg("emit file")
	return nil
}

func (wft *WriteFileToDiskTask) shouldSkip() bool {
	return fileExists(wft.getFullPath()) && !wft.overwrite
}

func (wft *WriteFileToDiskTask) getFullPath() string {
	return path.Join(wft.folder, wft.name+wft.suffix)
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
	log.Debug().Str("file", filePath).Int64("bytes", n).Msg("write file")
	return nil
}

func (wft *WriteFileToDiskTask) mkdirIfNecessary() error {
	dir := wft.folder
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

type ReadFromProviderTask struct {
	provider FileProvider
	out      io.Writer
}

func NewReadFromProviderTask(provider FileProvider, out io.Writer) *ReadFromProviderTask {
	return &ReadFromProviderTask{provider: provider, out: out}
}

func (t *ReadFromProviderTask) Run() error {
	fd, err := t.provider.Get()
	if err != nil {
		return fmt.Errorf("get file from provider: %w", err)
	}

	if _, err := io.Copy(t.out, fd); err != nil {
		return fmt.Errorf("copy file from provider: %w", err)
	}
	return nil
}
