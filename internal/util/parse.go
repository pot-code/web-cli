package util

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/pot-code/web-cli/internal/constant"
	"golang.org/x/mod/modfile"
)

var (
	ErrGoModNotFound = errors.New("can't locate go.mod")
)

type GoModMeta struct {
	Author      string
	ProjectName string
	Requires    []*modfile.Require
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

	mp, err := modfile.Parse(constant.GoModFile, content, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse mod file")
	}

	meta := new(GoModMeta)
	meta.Requires = mp.Require

	url := mp.Module.Mod.Path
	parts := strings.Split(url, "/")
	meta.Author = parts[len(parts)-2]
	meta.ProjectName = parts[len(parts)-1]

	if meta.Author == "" || meta.ProjectName == "" {
		return nil, errors.New("failed to extrace meta data from go.mod")
	}
	return meta, nil
}
