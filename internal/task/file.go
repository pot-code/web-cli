package task

import (
	"bytes"
	"io/fs"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/transformer"
	"github.com/pot-code/web-cli/internal/util"
	log "github.com/sirupsen/logrus"
)

type FileGenerator struct {
	Name         string
	Data         *bytes.Buffer
	Overwrite    bool
	transformers []transformer.Transformer
}

var _ Task = &FileGenerator{}

func (fg *FileGenerator) Run() error {
	if err := fg.validateRequest(); err != nil {
		return errors.Wrapf(err, "validate error")
	}

	if fg.shouldSkip() {
		log.WithFields(log.Fields{"file": fg.Name, "overwrite": fg.Overwrite}).Info("emit file [skipped]")
		return nil
	}

	if err := fg.applyTransformers(); err != nil {
		return err
	}

	if err := fg.writeToDisk(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"overwrite": fg.Overwrite,
		"file":      fg.Name,
	}).Infof("emit file")
	return nil
}

// UseTransformers the order matters
func (fg *FileGenerator) UseTransformers(tc ...transformer.Transformer) *FileGenerator {
	fg.transformers = append(fg.transformers, tc...)
	return fg
}

func (fg *FileGenerator) applyTransformers() error {
	trans := fg.transformers
	var tf transformer.TransformerFunc = func(result *bytes.Buffer) error {
		fg.Data = result
		return nil
	}
	tf = transformer.ApplyTransformers(tf, trans...)
	return errors.Wrap(tf(fg.Data), "failed to apply transformers")
}

func (fg *FileGenerator) shouldSkip() bool {
	return util.FileExists(fg.Name) && !fg.Overwrite
}

func (fg *FileGenerator) validateRequest() error {
	if fg.Name == "" {
		return errors.New("empty path")
	}
	if fg.Data.Len() == 0 {
		return errors.New("empty data")
	}
	return nil
}

func (fg *FileGenerator) writeToDisk() error {
	if err := fg.mkdirIfNecessary(); err != nil {
		return err
	}

	fn := fg.Name
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create target file")
	}

	w, err := fg.Data.WriteTo(fd)
	if err != nil {
		return errors.Wrapf(err, "failed to write data to '%s'", fn)
	}

	log.WithFields(log.Fields{"length": w, "file": fn}).Debug("write info")
	return nil
}

func (fg *FileGenerator) mkdirIfNecessary() error {
	dir := path.Dir(fg.Name)
	if dir == "" {
		return nil
	}
	if util.FileExists(dir) {
		return nil
	}
	return errors.Wrapf(os.MkdirAll(dir, fs.ModePerm), "failed to make '%s'", dir)
}

func BatchFileGenerationTask(gens []*FileGenerator, trans ...transformer.Transformer) []Task {
	tasks := make([]Task, len(gens))
	for i, g := range gens {
		tasks[i] = g.UseTransformers(trans...)
	}
	return tasks
}

type FileGenerationTree struct {
	name     string
	parent   *FileGenerationTree
	branches []*FileGenerationTree
	gens     []*FileGenerator
}

func NewFileGenerationTree(root string) *FileGenerationTree {
	return &FileGenerationTree{name: root}
}

func (ftr *FileGenerationTree) Branch(name string) *FileGenerationTree {
	nb := NewFileGenerationTree(name)
	nb.parent = ftr
	ftr.branches = append(ftr.branches, nb)
	return nb
}

func (ftr *FileGenerationTree) Up() *FileGenerationTree {
	return ftr.parent
}

func (ftr *FileGenerationTree) AddNodes(gens ...*FileGenerator) *FileGenerationTree {
	ftr.gens = append(ftr.gens, gens...)
	return ftr
}

func (ftr *FileGenerationTree) Flatten() []*FileGenerator {
	var result []*FileGenerator

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

func (ftr *FileGenerationTree) flatten(result *[]*FileGenerator, branch *FileGenerationTree, prefix string) {
	name := branch.name
	prefix = path.Join(prefix, name)
	for _, r := range branch.gens {
		r.Name = path.Join(prefix, r.Name)
		*result = append(*result, r)
	}

	for _, b := range branch.branches {
		ftr.flatten(result, b, prefix)
	}
}
