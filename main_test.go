package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// set to false for debugging
const suppressLog = true

// withProgram performs test setup and cleanup including: suppresses log output, updates PATH environment
// variable to preempt system restic with mock one, clears mock log, creates Program struct.
func withProgram(configFile string, test func(*Program)) {
	savedOut := log.Out
	savedPath := os.Getenv("PATH")

	defer func() {
		log.Out = savedOut
		_ = os.Setenv("PATH", savedPath)
	}()

	// suppress logs
	if suppressLog {
		log.Out = ioutil.Discard
	}

	// preempt path with mock scripts
	scriptsPath, err := filepath.Abs("scripts")
	if err == nil {
		err = os.Setenv("PATH", scriptsPath+string(os.PathListSeparator)+savedPath)
	}

	// clear restic command log
	_ = os.Remove("tmp/commands.log")

	cli := CLI{
		ConfigFile: configFile,
		StateFile: "tmp/state.json",
	}
	cli.Log.Level = "debug"

	p, err := NewProgram(&cli)

	if err != nil {
		panic(err)
	}

	test(p)
}

func resticCommands() []string {
	data, err := ioutil.ReadFile("tmp/commands.log")

	if err != nil {
		return []string{}
	} else {
		return strings.Split(strings.TrimSpace(string(data)), "\n")
	}
}

func TestRun(t *testing.T) {
	t.Run("Runs all commands", func(t *testing.T) {
		withProgram("test/simple_config.yml", func(p *Program) {
			cmd := RunCmd{}
			err := cmd.Run(p)

			assert.NoError(t, err)
			assert.Equal(t,
				[]string{
					"RESTIC_VAR1=rv1 RESTIC_VAR2=rv2 restic backup --f1 v1 --f2 a1 a2",
					"RESTIC_VAR1=rv1 RESTIC_VAR2=rv2 STDIN=Hello\\ World restic backup --f1 v1 --stdin",
				},
				resticCommands(),
			)
		})
	})

	t.Run("Runs specified command", func(t *testing.T) {
		withProgram("test/simple_config.yml", func(p *Program) {
			cmd := RunCmd{Commands: []string{"basic"}}
			err := cmd.Run(p)

			assert.NoError(t, err)
			assert.Equal(t,
				[]string{
					"RESTIC_VAR1=rv1 RESTIC_VAR2=rv2 restic backup --f1 v1 --f2 a1 a2",
				},
				resticCommands(),
			)
		})
	})

	t.Run("Does not run commands if dry run", func(t *testing.T) {
		withProgram("test/simple_config.yml", func(p *Program) {
			p.DryRun = true
			cmd := RunCmd{Commands: []string{"basic"}}
			err := cmd.Run(p)

			assert.NoError(t, err)
			assert.Empty(t, resticCommands())
		})

	})

	t.Run("Uses state file for commands with autoparent", func(t *testing.T) {
		err := ioutil.WriteFile(
			"tmp/state.json",
			[]byte("{\"snapshots\": {\"with_parent\": \"abc123\"}}"),
			0644,
		)
		assert.NoError(t, err)

		withProgram("test/autoparent_config.yml", func(p *Program) {
			cmd := RunCmd{}
			err = cmd.Run(p)

			assert.NoError(t, err)
			assert.Equal(t,
				[]string{
					"restic backup --parent abc123 with_parent",
					"restic backup no_parent",
				},
				resticCommands(),
			)
		})

		// ensure parent flags reflect snapshot IDs from previous run
		withProgram("test/autoparent_config.yml", func(p *Program) {
			cmd := RunCmd{}
			err = cmd.Run(p)

			assert.NoError(t, err)
			assert.Equal(t,
				[]string{
					"restic backup --parent bedabb1e with_parent",
					"restic backup --parent bedabb1e no_parent",
				},
				resticCommands(),
			)
		})
	})

	t.Run("Command failure returns error", func(t *testing.T) {
		withProgram("test/failure1_config.yml", func(p *Program) {
			cmd := RunCmd{}
			err := cmd.Run(p)

			assert.EqualError(t, err, "one or more commands failed: command_failure")
			assert.Equal(t, []string{"restic forget --fail"}, resticCommands())
		})

		withProgram("test/failure2_config.yml", func(p *Program) {
			cmd := RunCmd{}
			err := cmd.Run(p)

			assert.EqualError(t, err, "one or more commands failed: stdin_failure")
		})
	})
}
