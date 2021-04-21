package util

import (
	"os"
	"path"
	"strings"

	"github.com/pot-code/web-cli/pkg/core"
	log "github.com/sirupsen/logrus"
	tp "github.com/xlab/treeprint"
)

type TaskComposer struct {
	root     string // generation root
	files    []*core.FileDesc
	commands []*core.Command
	runner   core.Generator
	cleaned  bool
}

var _ core.Generator = &TaskComposer{}

func NewTaskComposer(root string, files ...*core.FileDesc) *TaskComposer {
	return &TaskComposer{root: root, files: files, cleaned: false}
}

// AddFile add file task
func (tc *TaskComposer) AddFile(fd *core.FileDesc) *TaskComposer {
	tc.files = append(tc.files, fd)
	return tc
}

// AddCommand add command
func (tc *TaskComposer) AddCommand(cmd *core.Command) *TaskComposer {
	tc.commands = append(tc.commands, cmd)
	return tc
}

func (tc *TaskComposer) Run() error {
	log.Debugf("generation tree:\n%s", tc.getGenerationTree())
	runner := tc.makeRunner()
	tc.runner = runner

	return runner.Run()
}

func (tc *TaskComposer) Cleanup() error {
	if tc.cleaned {
		return nil
	}
	if tc.runner == nil {
		return nil
	}
	if tc.root != "" {
		root := tc.root
		log.Debugf("removing folder '%s'", root)
		err := os.RemoveAll(root)
		if err != nil {
			log.WithFields(log.Fields{"error": err.Error(), "folder": root}).Debug("[cleanup]failed to cleanup")
		}
		return err
	}
	return tc.runner.Cleanup()
}

// MakeRunner make a task runner
func (tc *TaskComposer) makeRunner() core.Generator {
	var tasks []core.Executor

	for _, fd := range tc.files {
		if tc.root != "" {
			fd.Path = path.Join(tc.root, fd.Path)
		}
		tasks = append(tasks, core.NewFileGenerator(fd))
	}
	for _, cmd := range tc.commands {
		tasks = append(tasks, core.NewCmdExecutor(cmd))
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
