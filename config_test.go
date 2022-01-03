package main

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestCommandList(t *testing.T) {
	t.Run("UnmarshalYAML", func(t *testing.T) {
		commandList := CommandList{}

		err := yaml.UnmarshalStrict([]byte(`
a: {command: backup}
b: {command: forget}
c: {command: prune}
`), &commandList)

		assert.NoError(t, err)
		assert.Equal(t, commandList, CommandList{
			Command{Name: "a", Command: "backup"},
			Command{Name: "b", Command: "forget"},
			Command{Name: "c", Command: "prune"},
		})
	})
}
