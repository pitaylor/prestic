package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func (c *Command) Run(dryRun bool) (result CommandResult, err error) {
	resticCmd := c.ResticCmd()
	stdinCmd := c.StdinCmd()

	if stdinCmd != nil {
		log.WithFields(logrus.Fields{
			"command": stdinCmd,
			"env": stdinCmd.Env,
		}).Debug("Stdin command")
	}

	log.WithFields(logrus.Fields{
		"command": resticCmd,
		"env": resticCmd.Env,
	}).Debug("Restic command")

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
		logWriter := log.WriterLevel(logrus.DebugLevel)

		defer func(logWriter *io.PipeWriter) {
			if err := logWriter.Close(); err != nil {
				log.WithError(err).Error("Unable to close logger")
			}
		}(logWriter)

		teeReader := io.TeeReader(reader, logWriter)
		result.SnapshotId, err = c.snapshotScan(teeReader)

		if err != nil {
			log.WithError(err).Error("Problem scanning output")
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

	for _, flag := range c.Flags {
		cmd = append(cmd, flag.CommandArg()...)
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

func (f *Flag) CommandArg() []string {
	flag := f.Name

	if !strings.HasPrefix(flag, "-") {
		flag = "--" + flag
	}

	if boolVal, ok := f.Value.(bool); ok && boolVal {
		return []string{flag}
	} else {
		return []string{flag, os.ExpandEnv(fmt.Sprintf("%v", f.Value))}
	}
}
