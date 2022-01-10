package main

type EnvMap map[string]string
type CommandList []Command
type FlagList []Flag
type SnapshotMap map[string]string

// Command is a restic command and its arguments.
type Command struct {
	// Name is the unique command name specified in the configuration file.
	Name string

	// Command is the restic command, one of "backup", "forget", "prune", etc.
	Command string `yaml:""`

	// Env specifies the environment for the command.
	Env EnvMap `yaml:",omitempty"`

	// Flags specifies the command line flags. It is a map of the flag name to the
	// flag value. If the flag is not prefixed with a hyphen, a double hyphen is
	// assumed. The flag value is converted to its string representation except for the
	// boolean value true, which is not emitted.
	Flags FlagList `yaml:",omitempty"`

	// Args specifies the positional command line arguments.
	Args []string `yaml:",omitempty"`

	// Stdin specifies a program that is piped to the command's standard input.
	Stdin string `yaml:",omitempty"`

	// AutoParent specifies whether parent flag should be set automatically to snapshot ID from the last run.
	AutoParent bool `yaml:",omitempty"`
}

// Flag is a restic command flag.
type Flag struct {
	Name  string
	Value interface{}
}

// CommandResult is the result from running a restic command.
type CommandResult struct {
	// SnapshotId is the snapshot ID parsed from the command output.
	SnapshotId string
}

// Config is the program configuration.
type Config struct {
	// Commands specifies the restic commands to run.
	Commands CommandList `yaml:",omitempty"`
}

// State is program state that is persisted across program runs.
type State struct {
	// Snapshots is a map of command name to snapshot ID.
	Snapshots SnapshotMap `json:"snapshots,omitempty"`
}

type Program struct {
	Config      Config
	DryRun      bool
	StateFile   string
}
