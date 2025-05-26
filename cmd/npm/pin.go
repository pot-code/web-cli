package npm

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/pot-code/web-cli/pkg/command"
	"github.com/pot-code/web-cli/pkg/fileutil"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var PinPackageCmd = command.NewCommand("pin", "根据 pnpm-lock.yaml 锁定 package.json 依赖版本",
	any(nil),
	command.WithAlias([]string{"p"}),
).AddHandler(
	command.InlineHandler[any](
		func(c *cli.Context, config any) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get cwd: %w", err)
			}

			lockFilePath := path.Join(cwd, "pnpm-lock.yaml")
			if !fileutil.Exists(lockFilePath) {
				return fmt.Errorf("pnpm-lock.yaml not found")
			}

			lockFile, err := os.ReadFile(lockFilePath)
			if err != nil {
				return fmt.Errorf("read pnpm-lock.yaml: %w", err)
			}

			var lockFileSchema PnpmLockYamlSchema
			if err := yaml.Unmarshal(lockFile, &lockFileSchema); err != nil {
				return fmt.Errorf("unmarshal pnpm-lock.yaml: %w", err)
			}

			packageJsonFilePath := path.Join(cwd, "package.json")
			if !fileutil.Exists(packageJsonFilePath) {
				return fmt.Errorf("package.json not found")
			}

			packageJsonFile, err := os.ReadFile(packageJsonFilePath)
			if err != nil {
				return fmt.Errorf("read package.json: %w", err)
			}

			var packageJsonSchema PackageJsonSchema
			if err := json.Unmarshal(packageJsonFile, &packageJsonSchema); err != nil {
				return fmt.Errorf("unmarshal package.json: %w", err)
			}

			installedPackages := lockFileSchema.Importers["."].Dependencies
			for k := range packageJsonSchema.Dependencies {
				packageJsonSchema.Dependencies[k] = formatInstalledVersion(installedPackages[k].Version)
			}

			installedDevPackages := lockFileSchema.Importers["."].DevDependencies
			for k := range packageJsonSchema.DevDependencies {
				packageJsonSchema.DevDependencies[k] = formatInstalledVersion(installedDevPackages[k].Version)
			}

			packageJsonFile, err = json.MarshalIndent(packageJsonSchema, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal package.json: %w", err)
			}

			if err := os.WriteFile(packageJsonFilePath, packageJsonFile, fs.ModePerm); err != nil {
				return fmt.Errorf("write package.json: %w", err)
			}

			log.Info().Str("file", packageJsonFilePath).Msg("更新 package.json")
			return nil
		},
	),
).Create()

func formatInstalledVersion(v string) string {
	pi := strings.IndexByte(v, '(')
	if pi != -1 {
		v = v[:pi]
	}
	return v
}
