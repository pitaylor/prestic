package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (p *Program) GetState() *State {
	state := State{}

	if _, err := os.Stat(p.StateFile); os.IsNotExist(err) {
		log.WithFields(logrus.Fields{"file": p.StateFile}).Debug("State file not found")
	} else {
		data, err := ioutil.ReadFile(p.StateFile)
		if err == nil {
			err = json.Unmarshal(data, &state)
		}
		if err != nil {
			log.WithError(err).Error("Unable to read state file")
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
		dir := filepath.Dir(p.StateFile)

		src, err := os.Stat(dir)

		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
		} else if err == nil && !src.Mode().IsDir() {
			err = errors.New(fmt.Sprintf("state file path \"%v\" is not a directory", dir))
		}
	}

	if err == nil {
		err = ioutil.WriteFile(p.StateFile, data, 0600)
	}

	if err != nil {
		log.WithError(err).Error("Unable to update state file")
	}
}
