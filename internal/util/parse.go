package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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
			return nil, ErrGoModNotFound
		}
		return nil, fmt.Errorf("open go.mod: %w", err)
	}

	content, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("read go.mod: %w", err)
	}

	mp, err := modfile.Parse(GoModFile, content, nil)
	if err != nil {
		return nil, fmt.Errorf("parse mod file: %w", err)
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
