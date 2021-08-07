package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func (c *Command) CreateCmd(program string, preArgs ...string) *exec.Cmd {
	var args []string

	args = append(args, preArgs...)
	args = append(args, c.CommandArgs()...)

	cmd := exec.Command(program, args...)
	cmd.Env = make([]string, len(c.Env))

	for k, v := range c.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", k, os.ExpandEnv(v)))
	}

	return cmd
}

func (c *Command) CommandArgs() []string {
	var cmd []string

	for _, flag := range sortedKeys(c.Flags) {
		val := c.Flags[flag]

		if !strings.HasPrefix(flag, "-") {
			flag = "--" + flag
		}

		if boolVal, ok := val.(bool); ok && boolVal {
			cmd = append(cmd, flag)
		} else {
			cmd = append(cmd, flag, os.ExpandEnv(fmt.Sprintf("%v", val)))
		}
	}

	for _, arg := range c.Args  {
		cmd = append(cmd, os.ExpandEnv(arg))
	}

	return cmd
}

func sortedKeys(m map[string]interface{}) []string {
	sorted := make([]string, 0, len(m))

	for key := range m {
		sorted = append(sorted, key)
	}

	sort.Strings(sorted)

	return sorted
}
