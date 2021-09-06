package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func (p *Program) GetState() *State {
	state := State{}

	if _, err := os.Stat(p.StateFile); err != nil {
		log.Printf("No state file found: %v", p.StateFile)
	} else {
		data, err := ioutil.ReadFile(p.StateFile)
		if err == nil {
			err = json.Unmarshal(data, &state)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	if state.Snapshots == nil {
		state.Snapshots = make(SnapshotMap)
	}

	return &state
}

func (p *Program) UpdateState(updateFunc func(*State)) {
	state := p.GetState()

	updateFunc(state)

	data, err := json.Marshal(state)

	if err == nil {
		err = ioutil.WriteFile(p.StateFile, data, 0600)
	}

	if err != nil {
		log.Fatal(err)
	}
}
