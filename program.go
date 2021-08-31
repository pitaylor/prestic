package main

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func (p *Program) LoadConfig() {
	data, err := ioutil.ReadFile(p.ConfigFile)
	if err != nil {
		log.Panic(err)
	}

	err = yaml.UnmarshalStrict(data, &p.Config)
	if err != nil {
		log.Panic(err)
	}

	for _, action := range p.Config.Actions {
		action.Configure(&p.Config)
	}
}

func (p *Program) LoadState() {
	_, err := os.Stat(p.StateFile)
	if err != nil {
		log.Printf("No state file found: %v", p.StateFile)
	} else {
		data, err := ioutil.ReadFile(p.StateFile)
		if err == nil {
			err = json.Unmarshal(data, &p.State)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	if p.State.Snapshots == nil {
		p.State.Snapshots = make(SnapshotMap)
	}

	log.Printf("Loaded state: %v", p.State)
}

func (p *Program) SaveState() {
	data, err := json.Marshal(p.State)
	if err == nil {
		err = ioutil.WriteFile(p.StateFile, data, 0600)
	}

	if err != nil {
		log.Print(err)
		log.Printf("Error saving state: %v", p.StateFile)
		return
	}

	log.Printf("Saved state: %v", p.State)
}

func (p *Program) RunAll() {
	for _, action := range p.Config.Actions {
		if !p.DryRun {
			p.Run(action)
		} else {
			resticCmd := action.ResticCmd(p.resticFlags(action)...)
			stdinCmd := action.StdinCmd()
			if stdinCmd != nil {
				log.Printf("Stdin Command: %v, Env: %v", stdinCmd, stdinCmd.Env)
			}
			log.Printf("Restic Command: %v, Env: %v", resticCmd, resticCmd.Env)
		}
	}
}

func (p *Program) Run(action *Action) {
	result, err := action.Run(p.resticFlags(action)...)
	if err != nil {
		log.Printf("Command failed: %v", err)
	}

	if action.SnapshotKey != "" && result.SnapshotId != "" {
		p.State.Snapshots[action.SnapshotKey] = result.SnapshotId
		p.SaveState()
	}
}

func (p *Program) resticFlags(action *Action) []string {
	var extraArgs []string

	if action.SnapshotKey != "" {
		if snapshotId := p.State.Snapshots[action.SnapshotKey]; snapshotId != "" {
			extraArgs = append(extraArgs, "--parent", snapshotId)
		}
	}

	return extraArgs
}
