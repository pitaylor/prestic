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

func (a *Action) Configure(config *Config, snapshots SnapshotMap) {
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

	if snapshotId := snapshots[a.SnapshotKey]; snapshotId != "" && a.SnapshotKey != "" {
		a.Context.Flags["parent"] = snapshotId
	}
}

func (a *Action) Run(dryRun bool) (result ActionResult, err error) {
	resticCmd := a.ResticCmd()
	stdinCmd := a.StdinCmd()

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
		result.SnapshotId, err = a.snapshotScan(teeReader)

		if err != nil {
			log.Printf("Scan error: %v", err)
		}

		err = resticCmd.Wait()
	}

	return
}

func (a *Action) snapshotScan(reader io.Reader) (snapshotId string, err error) {
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

func (a *Action) ResticCmd() *exec.Cmd {
	var args []string
	args = append(args, a.Command)
	args = append(args, a.Context.CommandArgs()...)

	cmd := exec.Command("restic", args...)
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
