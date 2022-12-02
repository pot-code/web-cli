package util

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
)

const GoModFile = "go.mod"

var (
	ErrGoModNotFound = errors.New("can't locate go.mod")
)

type GoModMeta struct {
	author  string
	project string
	mf      *modfile.File
}

func ParseGoMod(path string) (*GoModMeta, error) {
	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithStack(ErrGoModNotFound)
		}
		return nil, errors.Wrap(err, "failed to open go.mod")
	}

	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	mp, err := modfile.Parse(GoModFile, content, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse mod file")
	}

	url := mp.Module.Mod.Path
	parts := strings.Split(url, "/")
	author := parts[len(parts)-2]
	project := parts[len(parts)-1]

	return &GoModMeta{
		mf:      mp,
		author:  author,
		project: project,
	}, nil
}

func (gm *GoModMeta) GetAuthor() string {
	return gm.author
}

func (gm *GoModMeta) GetProject() string {
	return gm.project
}

func (gm *GoModMeta) Contains(module string) bool {
	for _, r := range gm.mf.Require {
		if r.Mod.Path == module {
			return true
		}
	}
	return false
}
