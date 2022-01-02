package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func (p *Program) LoadConfig() {
	data, err := ioutil.ReadFile(p.ConfigFile)

	if err == nil {
		err = yaml.UnmarshalStrict(data, &p.Config)
	}

	if err != nil {
		log.Fatal(err)
	}

	snapshots := p.GetState().Snapshots

	for _, action := range p.Config.Actions {
		action.Configure(&p.Config, snapshots)
	}
}

func (p *Program) RunAll() {
	for _, action := range p.Config.Actions {
		p.Run(action)
	}
}

func (p *Program) Run(action *Action) {
	result, err := action.Run(p.DryRun)

	if err != nil {
		log.Printf("Command failed: %v", err)
	}

	if !p.DryRun && action.SnapshotKey != "" && result.SnapshotId != "" {
		p.UpdateState(func(state *State) {
			state.Snapshots[action.SnapshotKey] = result.SnapshotId
		})
	}
}
