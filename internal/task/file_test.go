package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTreeBranch_Flatten(t *testing.T) {
	tree := NewFileGenerationTree("root").
		Branch("sub").AddNodes(&FileGenerator{Name: "file"}).Up().
		Branch("sub1").AddNodes(&FileGenerator{Name: "file1"}).Up().
		Branch("sub2").
		AddNodes(
			&FileGenerator{Name: "file2"},
			&FileGenerator{Name: "file3"},
		)

	expect := []string{"root/sub/file", "root/sub1/file1", "root/sub2/file2", "root/sub2/file3"}
	actual := make([]string, 0)
	for _, r := range tree.Flatten() {
		actual = append(actual, r.Name)
	}

	assert.ElementsMatch(t, expect, actual)
}
