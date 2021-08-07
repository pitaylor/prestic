package main

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestCommand(t *testing.T) {
	t.Run("CommandArgs", func(t *testing.T) {
		config := Command{}
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
`), &config)

		assert.NoError(t, err)
		assert.Equal(
			t,
			[]string{"-f5", "5", "--f1", "1", "--f2", "--f3", "true", "--f4", "false", "a1", "a2"},
			config.CommandArgs(),
		)
	})

	t.Run("CommandArgs Empty", func(t *testing.T) {
		config := Command{}
		err := yaml.UnmarshalStrict([]byte(``), &config)

		assert.NoError(t, err)
		assert.Equal(t, []string(nil), config.CommandArgs())
	})
}
