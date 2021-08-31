package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func (a *Action) Configure(config *Config) {
	var contexts []*Context

	for _, preset := range strings.Split(a.Preset, ",") {
		preset = strings.TrimSpace(preset)

		if preset != "" {
			if context, ok := config.Presets[preset]; ok {
				contexts = append(contexts, context)
			} else {
				log.Fatalf("Invalid preset: %v", preset)
			}
		}
	}

	a.Context = *Merge(append(contexts, &a.Context))
}

func (a *Action) Run(resticFlags ...string) (result ActionResult, err error) {
	deferErr := func(e error) {
		if e != nil && err == nil {
			err = e
		}
	}

	resticCmd := a.ResticCmd(resticFlags...)

	if stdinCmd := a.StdinCmd(); stdinCmd != nil {
		resticCmd.Stdin, err = stdinCmd.StdoutPipe()
		if err == nil {
			err = stdinCmd.Start()
		}
		if err == nil {
			defer func() { deferErr(stdinCmd.Wait()) }()
		} else {
			return
		}
	}

	reader, err := resticCmd.StdoutPipe()
	if err == nil {
		err = resticCmd.Start()
	}
	if err == nil {
		defer func() { deferErr(resticCmd.Wait()) }()
		result.SnapshotId, err = a.snapshotScan(reader)
	}

	return
}

func (a *Action) snapshotScan(reader io.Reader) (snapshotId string, err error) {
	re := regexp.MustCompile(`snapshot ([0-9a-f]+) saved`)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		log.Println("RESTIC: " + scanner.Text())
		if match := re.FindSubmatch(scanner.Bytes()); match != nil {
			snapshotId = string(match[1])
		}
	}

	err = scanner.Err()
	return
}

func (a *Action) ResticCmd(flags ...string) *exec.Cmd {
	var args []string
	args = append(args, a.Command)
	args = append(args, flags...)
	args = append(args, a.Context.CommandArgs()...)

	// todo: fixme
	cmd := exec.Command("./restic.sh", args...)
	cmd.Env = a.Context.CommandEnv()
	cmd.Stderr = os.Stderr
	return cmd
}

func (a *Action) StdinCmd() *exec.Cmd {
	if a.Context.Stdin == "" {
		return nil
	}

	cmd := exec.Command("sh", "-c", a.Context.Stdin)
	cmd.Env = a.Context.CommandEnv()
	return cmd
}
