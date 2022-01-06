package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

func (p *Program) LoadConfig() {
	data, err := ioutil.ReadFile(p.ConfigFile)

	if err == nil {
		err = yaml.UnmarshalStrict(data, &p.Config)
	}

	if err != nil {
		log.WithError(err).Fatal("Unable to load config file")
	}
}

func (p *Program) ConfigureLogging() {
	if p.DryRun {
		log.Level = logrus.DebugLevel
	} else {
		level, err := logrus.ParseLevel(p.LogLevel)

		if err != nil {
			log.WithError(err).Fatal("Unable to parse log level")
		}

		log.Level = level
	}
}

func (p *Program) ConfigureParentFlags() {
	snapshots := p.GetState().Snapshots

	for i, command := range p.Config.Commands {
		if snapshotId := snapshots[command.Name]; snapshotId != "" && command.AutoParent {
			p.Config.Commands[i].Flags = append(command.Flags, Flag{Name: "parent", Value: snapshotId})
		}
	}
}

func (p *Program) RunAll() (errs []error) {
	for _, command := range p.Config.Commands {
		if err := p.Run(command); err != nil {
			errs = append(errs, err)
		}
	}

	log.WithFields(logrus.Fields{
		"success": len(p.Config.Commands) - len(errs),
		"failed":  len(errs),
	}).Info("Command summary")

	return
}

func (p *Program) Run(command Command) error {
	contextLog := log.WithFields(logrus.Fields{"name": command.Name})
	contextLog.Info("Command started")

	start := time.Now()
	result, err := command.Run(p.DryRun)

	if err != nil {
		contextLog.WithError(err).Error("Command failed")
	} else {
		contextLog.WithFields(logrus.Fields{
			"duration": time.Since(start).Nanoseconds(),
		}).Info("Command finished")
	}

	if !p.DryRun && command.AutoParent && result.SnapshotId != "" {
		p.UpdateState(func(state *State) {
			contextLog.WithFields(logrus.Fields{
				"snapshotId": result.SnapshotId,
			}).Debug("Updating state")
			state.Snapshots[command.Name] = result.SnapshotId
		})
	}

	return err
}
