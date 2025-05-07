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
	name      string // 文件名
	suffix    string // 文件后缀
	folder    string // 输出文件夹
	overwrite bool   // 是否覆盖
	data      io.Reader
}

type writeFileOption func(*WriteFileToDiskTask)

func WithFolder(folder string) writeFileOption {
	return func(w *WriteFileToDiskTask) {
		w.folder = folder
	}
}

func WithOverwrite(overwrite bool) writeFileOption {
	return func(w *WriteFileToDiskTask) {
		w.overwrite = overwrite
	}
}

func NewWriteFileToDiskTask(name string, suffix string, in io.Reader, options ...writeFileOption) *WriteFileToDiskTask {
	w := &WriteFileToDiskTask{
		name:   name,
		suffix: suffix,
		data:   in,
	}

	for _, opt := range options {
		opt(w)
	}
	return w
}

func (w *WriteFileToDiskTask) Run() error {
	if w.shouldSkip() {
		log.Info().Str("file", w.getFullPath()).Bool("overwrite", w.overwrite).Msg("emit file [skipped]")
		return nil
	}

	if err := w.write(); err != nil {
		return err
	}

	log.Info().Str("file", w.getFullPath()).Msg("emit file")
	return nil
}

func (w *WriteFileToDiskTask) shouldSkip() bool {
	return fileExists(w.getFullPath()) && !w.overwrite
}

func (w *WriteFileToDiskTask) getFullPath() string {
	return path.Join(w.folder, w.name+w.suffix)
}

func (w *WriteFileToDiskTask) write() error {
	if err := w.mkdirIfNecessary(); err != nil {
		return err
	}

	filePath := w.getFullPath()
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("open file [path: %s]: %w", filePath, err)
	}
	defer fd.Close()

	n, err := io.Copy(fd, w.data)
	if err != nil {
		return fmt.Errorf("write data to %s: %w", filePath, err)
	}
	log.Debug().Str("file", filePath).Int64("bytes", n).Msg("write file")
	return nil
}

func (w *WriteFileToDiskTask) mkdirIfNecessary() error {
	dir := w.folder
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
