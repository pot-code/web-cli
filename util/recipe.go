package util

import (
	"path"
	"strings"

	"github.com/pot-code/web-cli/core"
	tp "github.com/xlab/treeprint"
)

type GenerationMaterial struct {
	Path     string
	Provider core.DataProvider
}

func (gm GenerationMaterial) String() string {
	return gm.Path
}

type GenerationRecipe struct {
	root      string                // generation root
	materials []*GenerationMaterial // materials
}

func NewGenerationRecipe(root string, materials ...*GenerationMaterial) *GenerationRecipe {
	return &GenerationRecipe{root, materials}
}

// AddMaterial add material, chainable
func (gr *GenerationRecipe) AddMaterial(m *GenerationMaterial) *GenerationRecipe {
	gr.materials = append(gr.materials, m)
	return gr
}

// MakeGenerator map to a multi-task generator
func (gr *GenerationRecipe) MakeGenerator() core.Generator {
	tasks := make([]core.Generator, len(gr.materials))

	for i, m := range gr.materials {
		dst := m.Path
		if gr.root != "" {
			dst = path.Join(gr.root, dst)
		}
		tasks[i] = core.NewFileGenerator(dst, m.Provider)
	}
	return core.NewParallelGenerator(tasks...)
}

func (gr *GenerationRecipe) GetGenerationTree() tp.Tree {
	tree := tp.New()
	root := tree
	if gr.root != "" {
		root = tree.AddBranch(gr.root)
	}

	pathMap := make(map[string]tp.Tree)

	for _, m := range gr.materials {
		parts := strings.Split(m.Path, "/")
		gr.parseTree(parts, root, pathMap)
	}
	return tree
}

func (gr *GenerationRecipe) parseTree(parts []string, parent tp.Tree, lookup map[string]tp.Tree) {
	if len(parts) == 0 {
		return
	}

	val := parts[0]
	if len(parts) == 1 {
		parent.AddNode(val)
		return
	}

	if _, ok := lookup[val]; !ok {
		lookup[val] = parent.AddBranch(val)
	}
	gr.parseTree(parts[1:], lookup[val], lookup)
}
