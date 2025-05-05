package gomod

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
)

const GoModFile = "go.mod"

type GoModObject struct {
	mf *modfile.File
}

// Author 获取作者
func (gm *GoModObject) Author() string {
	modulePath := gm.mf.Module.Mod.Path
	if strings.Contains(modulePath, "github.com") {
		parts := strings.Split(modulePath, "/")
		return parts[1]
	}
	return ""
}

// ProjectName 获取项目名称
func (gm *GoModObject) ProjectName() string {
	modulePath := gm.mf.Module.Mod.Path
	if strings.Contains(modulePath, "github.com") {
		parts := strings.Split(modulePath, "/")
		return parts[2]
	}
	return modulePath
}

// HasModule 判断是否存在指定的依赖
func (gm *GoModObject) HasModule(module string) bool {
	for _, r := range gm.mf.Require {
		if r.Mod.Path == module {
			return true
		}
	}
	return false
}

var (
	ErrGoModFileNotFound = errors.New("找不到 go.mod 文件")
)

func ParseGoMod(path string) (*GoModObject, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrGoModFileNotFound
		}
		return nil, fmt.Errorf("read go.mod: %w", err)
	}

	mp, err := modfile.Parse(GoModFile, content, nil)
	if err != nil {
		return nil, fmt.Errorf("parse mod file: %w", err)
	}

	return &GoModObject{
		mf: mp,
	}, nil
}
