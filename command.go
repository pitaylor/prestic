package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func (c *Command) Run(dryRun bool) (result CommandResult, err error) {
	resticCmd := c.ResticCmd()
	stdinCmd := c.StdinCmd()

	if stdinCmd != nil {
		log.Printf("Stdin Command: %v, Env: %v", stdinCmd, stdinCmd.Env)
	}

	log.Printf("Restic Command: %v, Env: %v", resticCmd, resticCmd.Env)

	if dryRun {
		return
	}

	if stdinCmd != nil {
		var output []byte

		stdinCmd.Stderr = os.Stderr
		output, err = stdinCmd.Output()

		if err != nil {
			return
		}

		resticCmd.Stdin = strings.NewReader(string(output))
	}

	reader, err := resticCmd.StdoutPipe()

	if err == nil {
		err = resticCmd.Start()
	}

	if err == nil {
		teeReader := io.TeeReader(reader, os.Stdout)
		result.SnapshotId, err = c.snapshotScan(teeReader)

		if err != nil {
			log.Printf("Scan error: %v", err)
		}

		err = resticCmd.Wait()
	}

	return
}

func (c *Command) snapshotScan(reader io.Reader) (snapshotId string, err error) {
	re := regexp.MustCompile(`snapshot ([0-9a-f]+) saved`)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if match := re.FindSubmatch(scanner.Bytes()); match != nil {
			snapshotId = string(match[1])
		}
	}

	err = scanner.Err()
	return
}

func (c *Command) ResticCmd() *exec.Cmd {
	var args []string
	args = append(args, c.Command)
	args = append(args, c.CommandArgs()...)

	cmd := exec.Command("restic", args...)
	cmd.Env = c.CommandEnv()
	cmd.Stderr = os.Stderr
	return cmd
}

func (c *Command) StdinCmd() *exec.Cmd {
	if c.Stdin == "" {
		return nil
	}

	cmd := exec.Command("sh", "-c", c.Stdin)
	cmd.Env = c.CommandEnv()
	return cmd
}

func (c *Command) CommandArgs() []string {
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

func (c *Command) CommandEnv() []string {
	env := make([]string, len(c.Env)+1)

	// Include HOME so restic can find cache directory
	env = append(env, "HOME="+os.Getenv("HOME"))

	for k, v := range c.Env {
		env = append(env, fmt.Sprintf("%v=%v", k, os.ExpandEnv(v)))
	}

	return env
}
