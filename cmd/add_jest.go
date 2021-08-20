package cmd

import (
	"bytes"

	"github.com/pot-code/web-cli/pkg/core"
	"github.com/pot-code/web-cli/pkg/util"
	"github.com/pot-code/web-cli/templates"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const cmdJestName = "jest"

const eslintPatch = `
env: {
  'jest/globals': true,
},
plugins: ['jest'],
`

var addJestCmd = &cli.Command{
	Name:    cmdJestName,
	Usage:   "add Jest support",
	Aliases: []string{"j"},
	Action: func(c *cli.Context) error {
		cmd := addJestToReact()
		err := cmd.Run()
		if err != nil {
			cmd.Cleanup()
		} else {
			log.Warnf("add to .eslintrc.js: %s", eslintPatch)
		}
		return err
	},
}

func addJestToReact() core.Generator {
	return util.NewTaskComposer("",
		&core.FileDesc{
			Path: ".babelrc",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteJestBabelrc(&buf)
				return buf.Bytes()
			},
		},
		&core.FileDesc{
			Path: "jest.config.js",
			Data: func() []byte {
				var buf bytes.Buffer

				templates.WriteJestConfig(&buf)
				return buf.Bytes()
			},
		},
	).AddCommand(&core.Command{
		Bin: "npm",
		Args: []string{"add", "-D",
			"@babel/core",
			"babel-jest",
			"babel-loader",
			"@types/jest",
			"jest",
			"eslint-plugin-jest",
			"@testing-library/dom",
			"@testing-library/jest-dom",
			"@testing-library/react",
			"@testing-library/react-hooks",
			"@testing-library/user-event",
		},
	})
}
