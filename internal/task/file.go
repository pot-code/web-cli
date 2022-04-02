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

type FileRequest struct {
	Name      string
	Data      *bytes.Buffer
	Overwrite bool
}

type FileGenerator struct {
	fr           *FileRequest
	transformers []transformer.Transformer
}

func NewFileGenerator(req *FileRequest) *FileGenerator {
	return &FileGenerator{req, []transformer.Transformer{}}
}

func (fg *FileGenerator) Run() error {
	req := fg.fr
	if err := fg.validateRequest(); err != nil {
		return errors.Wrapf(err, "validate error")
	}

	if fg.shouldSkip() {
		log.WithFields(log.Fields{"file": req.Name, "overwrite": req.Overwrite}).Info("emit file [skipped]")
		return nil
	}

	if err := fg.applyTransformers(); err != nil {
		return err
	}

	if err := fg.writeToDisk(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"overwrite": req.Overwrite,
		"file":      req.Name,
	}).Infof("emit file")
	return nil
}

// UseTransformer the order matters
func (fg *FileGenerator) UseTransformer(tc ...transformer.Transformer) *FileGenerator {
	fg.transformers = append(fg.transformers, tc...)
	return fg
}

func (fg *FileGenerator) applyTransformers() error {
	req := fg.fr
	trans := fg.transformers
	var tf transformer.TransformerFunc = func(result *bytes.Buffer) error {
		req.Data = result
		return nil
	}
	tf = transformer.ApplyTransformers(tf, trans...)
	return errors.Wrap(tf(req.Data), "failed to apply transformers")
}

func (fg *FileGenerator) shouldSkip() bool {
	req := fg.fr
	return util.IsFileExist(req.Name) && !req.Overwrite
}

func (fg *FileGenerator) validateRequest() error {
	req := fg.fr
	if req.Name == "" {
		return errors.New("empty path")
	}
	if req.Data.Len() == 0 {
		return errors.New("empty data")
	}
	return nil
}

func (fg *FileGenerator) writeToDisk() error {
	req := fg.fr
	if err := fg.mkdirIfNecessary(); err != nil {
		return err
	}

	fd, err := os.OpenFile(req.Name, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create target file")
	}

	w, err := req.Data.WriteTo(fd)
	if err != nil {
		return errors.Wrapf(err, "failed to write data to '%s'", req.Name)
	}

	log.WithFields(log.Fields{"length": w, "file": req.Name}).Debug("write info")
	return nil
}

func (fg *FileGenerator) mkdirIfNecessary() error {
	dir := path.Dir(fg.fr.Name)
	if dir == "" {
		return nil
	}
	if util.IsFileExist(dir) {
		return nil
	}
	return errors.Wrapf(os.MkdirAll(dir, fs.ModePerm), "failed to make '%s'", dir)
}

func BatchFileTask(requests []*FileRequest, trans ...transformer.Transformer) []Task {
	tasks := make([]Task, len(requests))

	for i, r := range requests {
		tasks[i] = NewFileGenerator(r).UseTransformer(trans...)
	}
	return tasks
}

type FileRequestTree struct {
	name     string
	parent   *FileRequestTree
	branches []*FileRequestTree
	requests []*FileRequest
}

func NewFileRequestTree(root string) *FileRequestTree {
	return &FileRequestTree{name: root}
}

func (ftr *FileRequestTree) Branch(name string) *FileRequestTree {
	nb := NewFileRequestTree(name)
	nb.parent = ftr
	ftr.branches = append(ftr.branches, nb)
	return nb
}

func (ftr *FileRequestTree) Up() *FileRequestTree {
	return ftr.parent
}

func (ftr *FileRequestTree) AddNode(req ...*FileRequest) *FileRequestTree {
	ftr.requests = append(ftr.requests, req...)
	return ftr
}

func (ftr *FileRequestTree) Flatten() []*FileRequest {
	var result []*FileRequest

	root := ftr.root()
	ftr.flatten(&result, root, "")
	return result
}

func (ftr *FileRequestTree) root() *FileRequestTree {
	node := ftr
	for node.parent != nil {
		node = node.parent
	}
	return node
}

func (ftr *FileRequestTree) flatten(result *[]*FileRequest, branch *FileRequestTree, prefix string) {
	name := branch.name
	prefix = path.Join(prefix, name)
	for _, r := range branch.requests {
		r.Name = path.Join(prefix, r.Name)
		*result = append(*result, r)
	}

	for _, b := range branch.branches {
		ftr.flatten(result, b, prefix)
	}
}
