package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

func NewProgram(cli *CLI) (*Program, error) {
	p := Program{StateFile: cli.StateFile}

	err := p.LoadConfig(cli.ConfigFile)

	if err == nil {
		err = p.ConfigureLogging(cli.Log.Level)
	}

	if err == nil {
		p.ConfigureParentFlags()
	}

	return &p, err
}

func (p *Program) LoadConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)

	if err == nil {
		err = yaml.UnmarshalStrict(data, &p.Config)
	}

	return err
}

func (p *Program) ConfigureLogging(logLevel string) error {
	var err error
	if p.DryRun {
		log.Level = logrus.DebugLevel
	} else {
		var level logrus.Level
		if level, err = logrus.ParseLevel(logLevel); err == nil {
			log.Level = level
		}
	}

	return err
}

func (p *Program) ConfigureParentFlags() {
	snapshots := p.GetState().Snapshots

	for i, command := range p.Config.Commands {
		if snapshotId := snapshots[command.Name]; snapshotId != "" && command.AutoParent {
			p.Config.Commands[i].Flags = append(command.Flags, Flag{Name: "parent", Value: snapshotId})
		}
	}
}

func (p *Program) Run(command Command) error {
	contextLog := log.WithFields(logrus.Fields{"name": command.Name})
	contextLog.Info("Command started")

	start := time.Now()
	result, err := command.Run(p.DryRun)

	if err != nil {
		contextLog.WithError(err).Error("Command failed")
	} else {
		contextLog.WithField("duration", time.Since(start).Milliseconds()).Info("Command finished")
	}

	if !p.DryRun && command.AutoParent && result.SnapshotId != "" {
		p.UpdateState(func(state *State) {
			contextLog.WithField("snapshotId", result.SnapshotId).Debug("Updating state")
			state.Snapshots[command.Name] = result.SnapshotId
		})
	}

	return err
}
