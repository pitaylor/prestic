package main

import (
	"errors"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

var log = logrus.New()

func main() {
	err := wrappedMain(os.Args[1:]...)
	if err == flag.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		os.Exit(1)
	}
}

func wrappedMain(args ...string) error {
	p := Program{}

	presticDir := ".prestic"

	if homeDir, err := os.UserHomeDir(); err == nil {
		presticDir = path.Join(homeDir, presticDir)
	}

	cli := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cli.StringVar(&p.ConfigFile, "config", path.Join(presticDir, "config.yml"), "Config file")
	cli.StringVar(&p.StateFile, "state", path.Join(presticDir, "state.json"), "State file")
	cli.StringVar(&p.LogLevel, "log-level", "info", "Log level: debug, error, warn, info (default info)")
	cli.BoolVar(&p.DryRun, "dry-run", false, "Perform a dry run")
	err := cli.Parse(args)

	if err != nil {
		return err
	}

	p.LoadConfig()
	p.ConfigureLogging()
	p.ConfigureParentFlags()

	if errs := p.RunAll(); len(errs) != 0 {
		return errors.New("prestic: one or more commands failed")
	}

	return nil
}