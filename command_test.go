package main

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestCommand(t *testing.T) {
	t.Run("CommandArgs", func(t *testing.T) {
		command := Command{}
		err := yaml.UnmarshalStrict([]byte(`
flags:
  f1: 1
  f2: true
  f3: "true"
  f4: false
  -f5: "5"
args:
  - a1
  - a2
`), &command)

		assert.NoError(t, err)
		assert.ElementsMatch(
			t,
			[]string{"--f1", "1", "--f2", "--f3", "true", "--f4", "false", "-f5", "5", "a1", "a2"},
			command.CommandArgs(),
		)
	})

	t.Run("CommandArgs Empty", func(t *testing.T) {
		command := Command{}
		err := yaml.UnmarshalStrict([]byte(``), &command)

		assert.NoError(t, err)
		assert.Equal(t, []string(nil), command.CommandArgs())
	})
}
