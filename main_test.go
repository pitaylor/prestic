package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrappedMain(t *testing.T) {
	log.Out = ioutil.Discard

	savedPath := os.Getenv("PATH")

	defer func() { assert.NoError(t, os.Setenv("PATH", savedPath)) }()

	scriptsPath, err := filepath.Abs("scripts")

	if err == nil {
		err = os.Setenv("PATH", scriptsPath+string(os.PathListSeparator)+savedPath)
	}

	assert.NoError(t, err)

	clearCommands := func() {
		_ = os.Remove("tmp/commands.log")
	}

	readCommands := func() []string {
		data, err := ioutil.ReadFile("tmp/commands.log")
		assert.NoError(t, err)

		return strings.Split(strings.TrimSpace(string(data)), "\n")
	}

	t.Run("Runs commands", func(t *testing.T) {
		clearCommands()

		err := wrappedMain("-config", "test/simple_config.yml")
		assert.NoError(t, err)
		assert.Equal(t,
			[]string{
				"RESTIC_VAR1=rv1 RESTIC_VAR2=rv2 restic backup --f1 v1 --f2 a1 a2",
				"RESTIC_VAR1=rv1 RESTIC_VAR2=rv2 STDIN=Hello\\ World restic backup --f1 v1 --stdin",
			},
			readCommands(),
		)
	})

	t.Run("Uses state file for commands with autoparent", func(t *testing.T) {
		clearCommands()

		err := ioutil.WriteFile(
			"tmp/state.json",
			[]byte("{\"snapshots\": {\"with_parent\": \"abc123\"}}"),
			0644,
		)
		assert.NoError(t, err)

		err = wrappedMain("-config", "test/autoparent_config.yml", "-state", "tmp/state.json")
		assert.NoError(t, err)

		assert.Equal(t,
			[]string{
				"restic backup --parent abc123 with_parent",
				"restic backup no_parent",
			},
			readCommands(),
		)

		clearCommands()

		err = wrappedMain("-config", "test/autoparent_config.yml", "-state", "tmp/state.json")
		assert.NoError(t, err)

		assert.Equal(t,
			[]string{
				"restic backup --parent bedabb1e with_parent",
				"restic backup --parent bedabb1e no_parent",
			},
			readCommands(),
		)
	})

	t.Run("Command failures cause program error", func(t *testing.T) {
		clearCommands()

		err = wrappedMain("-config", "test/failure1_config.yml", "-log-level", "debug")
		assert.Error(t, err, "prestic: one or more commands failed")

		err = wrappedMain("-config", "test/failure2_config.yml", "-log-level", "debug")
		assert.Error(t, err, "prestic: one or more commands failed")

		assert.Equal(t, []string{"restic forget --fail"}, readCommands())
	})
}
