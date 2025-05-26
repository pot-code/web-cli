package npm

type PackageJsonSchema struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Private         bool              `json:"private"`
	Type            string            `json:"type"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

type PnpmLockYamlSchema struct {
	LockFileVersion string `yaml:"lockfileVersion"`
	Settings        struct {
		AutoInstallPeers         bool `yaml:"autoInstallPeers"`
		ExcludeLinksFromLockFile bool `yaml:"excludeLinksFromLockFile"`
	}
	Importers map[string]struct {
		Dependencies map[string]struct {
			Specifier string `yaml:"specifier"`
			Version   string `yaml:"version"`
		} `yaml:"dependencies"`
		DevDependencies map[string]struct {
			Specifier string `yaml:"specifier"`
			Version   string `yaml:"version"`
		} `yaml:"devDependencies"`
	}
}
