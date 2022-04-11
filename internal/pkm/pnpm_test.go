package pkm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pnpm_Create(t *testing.T) {
	pm := newPnpm("pnpm")
	cmd := pm.Create("next-app", "test", []string{"--typescript"})

	assert.Equal(t, "pnpm create next-app -- --typescript test", cmd.String())
}

func Test_pnpm_Install(t *testing.T) {
	pm := newPnpm("pnpm")
	cmd := pm.Install([]string{"react"})

	assert.Equal(t, "pnpm add react", cmd.String())
}

func Test_pnpm_InstallDev(t *testing.T) {
	pm := newPnpm("pnpm")
	cmd := pm.InstallDev([]string{"react"})

	assert.Equal(t, "pnpm add -D react", cmd.String())
}

func Test_pnpm_Uninstall(t *testing.T) {
	pm := newPnpm("pnpm")
	cmd := pm.Uninstall([]string{"react"})

	assert.Equal(t, "pnpm rm react", cmd.String())
}
