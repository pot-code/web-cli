package util

import (
	"path"
	"strings"

	"github.com/pot-code/web-cli/pkg/core"
	log "github.com/sirupsen/logrus"
	tp "github.com/xlab/treeprint"
)

type TaskComposer struct {
	root     string // generation root
	files    []*core.FileDesc
	commands []*core.ShellCommand
}

var _ core.Runner = &TaskComposer{}

func NewTaskComposer(root string) *TaskComposer {
	return &TaskComposer{root: root}
}

// AddFile add file task
func (tc *TaskComposer) AddFile(fds ...*core.FileDesc) *TaskComposer {
	tc.files = append(tc.files, fds...)
	return tc
}

// AddCommand add command
func (tc *TaskComposer) AddCommand(cmds ...*core.ShellCommand) *TaskComposer {
	tc.commands = append(tc.commands, cmds...)
	return tc
}

func (tc *TaskComposer) Run() error {
	log.Debugf("generation tree:\n%s", tc.getGenerationTree())
	runner := tc.makeRunner()
	return runner.Run()
}

// MakeRunner make a task runner
func (tc *TaskComposer) makeRunner() core.Runner {
	var tasks []core.Runner

	for _, fd := range tc.files {
		if tc.root != "" {
			fd.Path = path.Join(tc.root, fd.Path)
		}
		tasks = append(tasks, core.NewFileGenerator(fd))
	}
	for _, cmd := range tc.commands {
		tasks = append(tasks, core.NewShellCmdExecutor(cmd))
	}
	return core.NewParallelRunner(tasks...)
}

func (tc *TaskComposer) getGenerationTree() tp.Tree {
	tree := tp.New()
	root := tree
	if tc.root != "" {
		root = tree.AddBranch(tc.root)
	}

	pathMap := make(map[string]tp.Tree)

	for _, m := range tc.files {
		parts := strings.Split(m.Path, "/")
		tc.parseTree(parts, root, pathMap)
	}
	return tree
}

func (tc *TaskComposer) parseTree(parts []string, parent tp.Tree, lookup map[string]tp.Tree) {
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
	tc.parseTree(parts[1:], lookup[val], lookup)
}
