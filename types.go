package main

type EnvMap map[string]string
type FlagMap map[string]interface{}

type Command struct {
	Env   EnvMap   `yaml:",omitempty"`
	Flags FlagMap  `yaml:",omitempty"`
	Args  []string `yaml:",omitempty"`
	Stdin string   `yaml:",omitempty"`
}

type Config struct {
	Presets map[string]Command `yaml:",omitempty"`
	Backups []Command          `yaml:",omitempty"`
	Forgets []Command          `yaml:",omitempty"`
	Prune   *Command           `yaml:",omitempty"`
}

type Program struct {
	Config     Config
	ConfigFile string
	DryRun     bool
}
