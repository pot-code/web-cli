package pkm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_yarn_Create(t *testing.T) {
	pm := newYarn("yarn")
	cmd := pm.Create("next-app", "test", []string{"--typescript"})

	assert.Equal(t, "yarn create next-app --typescript test", cmd.String())
}

func Test_yarn_Install(t *testing.T) {
	pm := newYarn("yarn")
	cmd := pm.Install([]string{"react"})

	assert.Equal(t, "yarn add react", cmd.String())
}

func Test_yarn_InstallDev(t *testing.T) {
	pm := newYarn("yarn")
	cmd := pm.InstallDev([]string{"react"})

	assert.Equal(t, "yarn add -D react", cmd.String())
}

func Test_yarn_Uninstall(t *testing.T) {
	pm := newYarn("yarn")
	cmd := pm.Uninstall([]string{"react"})

	assert.Equal(t, "yarn remove react", cmd.String())
}
