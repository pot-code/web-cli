package util

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type GoModMeta struct {
	Author, ProjectName string
}

func ParseGoMod(path string) (*GoModMeta, error) {
	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithMessage(err, "can't locate go.mod")
		}
		return nil, errors.Wrap(err, "failed to open go.mod")
	}

	var author, projectName string
	reader := bufio.NewScanner(fd)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "module") {
			url := strings.TrimPrefix(line, "module ")
			parts := strings.Split(url, "/")
			author = parts[len(parts)-2]
			projectName = parts[len(parts)-1]
			break
		}
	}
	if author == "" || projectName == "" {
		return nil, errors.New("failed to extrace meta data from go.mod")
	}
	return &GoModMeta{author, projectName}, nil
}
