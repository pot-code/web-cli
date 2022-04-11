package pkm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_npm_Install(t *testing.T) {
	pm := newNpm("npm")
	cmd := pm.Install([]string{"react"})

	assert.Equal(t, "npm install react", cmd.String())
}

func Test_npm_InstallDev(t *testing.T) {
	pm := newNpm("npm")
	cmd := pm.InstallDev([]string{"react"})

	assert.Equal(t, "npm install -D react", cmd.String())
}

func Test_npm_Uninstall(t *testing.T) {
	pm := newNpm("npm")
	cmd := pm.Uninstall([]string{"react"})

	assert.Equal(t, "npm rm react", cmd.String())
}

func Test_npm_Create(t *testing.T) {
	pm := newNpm("npm")
	cmd := pm.Create("next-app", "test", []string{"--typescript"})

	assert.Equal(t, "npx next-app --typescript test", cmd.String())
}
