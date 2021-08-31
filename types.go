package main

type EnvMap map[string]string
type FlagMap map[string]interface{}
type SnapshotMap map[string]string

// Context is an execution context for a restic command.
type Context struct {
	// Env specifies the environment for the command.
	Env EnvMap `yaml:",omitempty"`

	// Flags specifies the command line flags. It is a map of the flag name to the
	// flag value. If the flag is not prefixed with a hyphen, a double hyphen is
	// assumed. The flag value is converted to its string representation except for the
	// boolean value true, which is not emitted.
	Flags FlagMap `yaml:",omitempty"`

	// Args specifies the positional command line arguments.
	Args []string `yaml:",omitempty"`

	// Stdin specifies a program that is piped to the command's standard input.
	Stdin string `yaml:",omitempty"`
}

// Action is a restic command.
type Action struct {
	Command     string  `yaml:""`
	Preset      string  `yaml:",omitempty"`
	Context     Context `yaml:",inline"`
	SnapshotKey string  `yaml:"snapshot-key,omitempty"`
}

type Config struct {
	Presets map[string]*Context `yaml:",omitempty"`
	Actions []*Action           `yaml:",omitempty"`
}

type ActionResult struct {
	SnapshotId string
}

type State struct {
	Snapshots SnapshotMap `json:"snapshots,omitempty"`
}

type Program struct {
	Config     Config
	ConfigFile string
	DryRun     bool
	State      State
	StateFile  string
}
