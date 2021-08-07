package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func (p *Program) Load() {
	data, err := ioutil.ReadFile(p.ConfigFile)

	if err != nil {
		log.Panic(err)
	}

	err = yaml.UnmarshalStrict(data, &p.Config)

	if err != nil {
		log.Panic(err)
	}
}

func (p *Program) Run() {
	for _, backup := range p.Config.Backups {
		p.runCmd(backup.CreateCmd("restic", "backup", "--quiet"))
	}

	for _, forget := range p.Config.Forgets {
		p.runCmd(forget.CreateCmd("restic", "forget", "--quiet"))
	}

	if p.Config.Prune != nil {
		p.runCmd(p.Config.Prune.CreateCmd("restic", "prune", "--quiet"))
	}
}

func (p *Program) runCmd(cmd *exec.Cmd) {
	// Pass HOME so restic can find cache directory
	cmd.Env = append(cmd.Env, "HOME="+os.Getenv("HOME"))

	log.Printf("Command: %v, Env: %v", cmd, cmd.Env)

	if p.DryRun {
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Printf("ERROR: command failed! %v", err)
	}
}
