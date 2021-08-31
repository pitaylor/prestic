package main

import (
	"fmt"
	"os"
	"strings"
)

func (c *Context) CommandArgs() []string {
	var cmd []string

	for flag, val := range c.Flags {
		if !strings.HasPrefix(flag, "-") {
			flag = "--" + flag
		}

		if boolVal, ok := val.(bool); ok && boolVal {
			cmd = append(cmd, flag)
		} else {
			cmd = append(cmd, flag, os.ExpandEnv(fmt.Sprintf("%v", val)))
		}
	}

	for _, arg := range c.Args {
		cmd = append(cmd, os.ExpandEnv(arg))
	}

	return cmd
}

func (c *Context) CommandEnv() []string {
	env := make([]string, len(c.Env)+1)

	// Include HOME so restic can find cache directory
	env = append(env, "HOME="+os.Getenv("HOME"))

	for k, v := range c.Env {
		env = append(env, fmt.Sprintf("%v=%v", k, os.ExpandEnv(v)))
	}

	return env
}

func Merge(contexts []*Context) *Context {
	result := Context{}
	result.Env = make(EnvMap)
	result.Flags = make(FlagMap)

	for _, c := range contexts {
		for k, v := range c.Env {
			result.Env[k] = v
		}

		for k, v := range c.Flags {
			result.Flags[k] = v
		}

		result.Args = append(result.Args, c.Args...)
		result.Stdin = c.Stdin
	}

	return &result
}
