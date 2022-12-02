package task

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/util"
	log "github.com/sirupsen/logrus"
)

type TemplateProvider interface {
	Get() (io.ReadCloser, error)
}

var _ Task = (*TemplateRenderTask)(nil)

type TemplateRenderTask struct {
	Name      string
	Suffix    string
	Folder    string
	Provider  TemplateProvider
	Overwrite bool
	Pipelines []PipelineFn
	Data      interface{}
}

func NewTemplateRenderTask(name string, suffix string, folder string, provider TemplateProvider, overwrite bool, data interface{}) *TemplateRenderTask {
	return &TemplateRenderTask{
		Name:      name,
		Suffix:    suffix,
		Folder:    folder,
		Overwrite: overwrite,
		Provider:  provider,
		Data:      data,
	}
}

func (fg *TemplateRenderTask) Run() error {
	log.WithFields(log.Fields{
		"name":           fg.Name,
		"suffix":         fg.Suffix,
		"folder":         fg.Folder,
		"overwrite":      fg.Overwrite,
		"pipeline_count": len(fg.Pipelines),
	}).Debug("generating file")
	if err := fg.validateRequest(); err != nil {
		return err
	}

	if fg.shouldSkip() {
		log.WithFields(log.Fields{
			"file":      fg.getFullPath(),
			"overwrite": fg.Overwrite,
		}).Info("emit file [skipped]")
		return nil
	}

	b := new(bytes.Buffer)
	if err := fg.renderTemplate(b); err != nil {
		return err
	}

	d := new(bytes.Buffer)
	if err := fg.applyPipeline(b, d); err != nil {
		return err
	}

	if err := fg.writeToDisk(d); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"file": fg.getFullPath(),
	}).Info("emit file")
	return nil
}

// UseTransformers the order matters
func (fg *TemplateRenderTask) AddPipelineFn(pn PipelineFn) *TemplateRenderTask {
	fg.Pipelines = append(fg.Pipelines, pn)
	return fg
}

func (fg *TemplateRenderTask) shouldSkip() bool {
	return util.Exists(fg.getFullPath()) && !fg.Overwrite
}

func (fg *TemplateRenderTask) validateRequest() error {
	if fg.Name == "" {
		return errors.New("empty name")
	}
	return nil
}

func (fg *TemplateRenderTask) getFullPath() string {
	return path.Join(fg.Folder, fg.Name+fg.Suffix)
}

func (fg *TemplateRenderTask) writeToDisk(b *bytes.Buffer) error {
	if err := fg.mkdirIfNecessary(); err != nil {
		return err
	}

	filePath := fg.getFullPath()
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("open file [path: %s]: %w", filePath, err)
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(b)
	if err != nil {
		return fmt.Errorf("read from buffer: %w", err)
	}

	n, err := fd.Write(data)
	if err != nil {
		return errors.Wrapf(err, "failed to write data to '%s'", filePath)
	}
	log.WithFields(log.Fields{"bytes": n, "file": filePath}).Debug("write file")
	return nil
}

func (fg *TemplateRenderTask) mkdirIfNecessary() error {
	dir := fg.Folder
	if dir == "" {
		return nil
	}
	if util.Exists(dir) {
		return nil
	}
	return errors.Wrapf(os.MkdirAll(dir, fs.ModePerm), "failed to make '%s'", dir)
}

func (fg *TemplateRenderTask) applyPipeline(src io.Reader, dest io.Writer) error {
	p := NewPipeline()

	for _, pn := range fg.Pipelines {
		p.AddPipelineFn(pn)
	}
	p.AddPipelineFn(func(src io.Reader, dest io.Writer) error {
		r, err := ioutil.ReadAll(src)
		if err != nil {
			return fmt.Errorf("read all from src: %w", err)
		}

		_, err = dest.Write(r)
		if err != nil {
			return fmt.Errorf("write to dest: %w", err)
		}
		return nil
	})

	if err := p.Apply(src, dest); err != nil {
		return fmt.Errorf("apply pipeline: %w", err)
	}
	return nil
}

func (fg *TemplateRenderTask) renderTemplate(out io.Writer) error {
	p := fg.Provider

	fd, err := p.Get()
	if err != nil {
		return fmt.Errorf("get template: %w", err)
	}

	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return fmt.Errorf("read template data: %w", err)
	}
	defer fd.Close()

	r := NewDefaultTemplateRenderer()
	log.WithField("template", string(b)).Debug("render template")
	err = r.Render(&RenderRequest{
		Name:     fg.Name,
		Template: string(b),
		Data:     fg.Data,
	}, out)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	return nil
}

type FileGenerationTree struct {
	folder   string
	parent   *FileGenerationTree
	branches []*FileGenerationTree
	gens     []*TemplateRenderTask
}

func NewFileGenerationTree(root string) *FileGenerationTree {
	return &FileGenerationTree{folder: root}
}

func (ftr *FileGenerationTree) Branch(folder string) *FileGenerationTree {
	nb := NewFileGenerationTree(folder)
	nb.parent = ftr
	ftr.branches = append(ftr.branches, nb)
	return nb
}

func (ftr *FileGenerationTree) Up() *FileGenerationTree {
	return ftr.parent
}

func (ftr *FileGenerationTree) AddNodes(gens ...*TemplateRenderTask) *FileGenerationTree {
	ftr.gens = append(ftr.gens, gens...)
	return ftr
}

func (ftr *FileGenerationTree) Flatten() []*TemplateRenderTask {
	var result []*TemplateRenderTask

	root := ftr.root()
	ftr.flatten(&result, root, "")
	return result
}

func (ftr *FileGenerationTree) root() *FileGenerationTree {
	node := ftr
	for node.parent != nil {
		node = node.parent
	}
	return node
}

func (ftr *FileGenerationTree) flatten(result *[]*TemplateRenderTask, branch *FileGenerationTree, prefix string) {
	folder := branch.folder
	prefix = path.Join(prefix, folder)
	for _, r := range branch.gens {
		r.Name = path.Join(prefix, r.Name)
		*result = append(*result, r)
	}

	for _, b := range branch.branches {
		ftr.flatten(result, b, prefix)
	}
}
