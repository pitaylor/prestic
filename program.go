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

	for _, command := range p.Config.Commands {
		if snapshotId := snapshots[command.Name]; snapshotId != "" && command.AutoParent {
			command.Flags["parent"] = snapshotId
		}
	}
}

func (p *Program) RunAll() (errs []error) {
	for _, command := range p.Config.Commands {
		if err := p.Run(command); err != nil {
			errs = append(errs, err)
		}
	}

	return
}

func (p *Program) Run(command Command) error {
	result, err := command.Run(p.DryRun)

	if err != nil {
		log.Printf("Command failed: %v", err)
	}

	if !p.DryRun && command.AutoParent && result.SnapshotId != "" {
		p.UpdateState(func(state *State) {
			state.Snapshots[command.Name] = result.SnapshotId
		})
	}

	return err
}
