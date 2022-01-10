package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func (p *Program) GetState() *State {
	state := State{}

	if _, err := os.Stat(p.StateFile); err != nil {
		log.WithFields(logrus.Fields{"file": p.StateFile}).Debug("State file not found")
	} else {
		data, err := ioutil.ReadFile(p.StateFile)
		if err == nil {
			err = json.Unmarshal(data, &state)
		}
		if err != nil {
			log.WithError(err).Fatal("Unable to read state file")
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
		log.WithError(err).Fatal("Unable to update state file")
	}
}
