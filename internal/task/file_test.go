package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTreeBranch_Flatten(t *testing.T) {
	tree := NewFileGenerationTree("root").
		Branch("sub").AddNodes(&FileGenerationTask{Name: "file"}).Up().
		Branch("sub1").AddNodes(&FileGenerationTask{Name: "file1"}).Up().
		Branch("sub2").Branch("sub3").
		AddNodes(
			&FileGenerationTask{Name: "file2"},
			&FileGenerationTask{Name: "file3"},
		)

	expect := []string{"root/sub/file", "root/sub1/file1", "root/sub2/sub3/file2", "root/sub2/sub3/file3"}
	actual := make([]string, 0)
	for _, r := range tree.Flatten() {
		actual = append(actual, r.Name)
	}

	assert.ElementsMatch(t, expect, actual)
}
