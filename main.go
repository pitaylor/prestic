package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

var log = logrus.New()

func main() {
	p := Program{}

	presticDir := ".prestic"

	if homeDir, err := os.UserHomeDir(); err == nil {
		presticDir = path.Join(homeDir, presticDir)
	}

	flag.StringVar(&p.ConfigFile, "config", path.Join(presticDir, "config.yml"), "Config file")
	flag.StringVar(&p.StateFile, "state", path.Join(presticDir, "state.json"), "State file")
	flag.StringVar(&p.LogLevel, "log-level", "info", "Log level: debug, error, warn, info (default info)")
	flag.BoolVar(&p.DryRun, "dry-run", false, "Perform a dry run")
	flag.Parse()

	p.LoadConfig()
	p.ConfigureLogging()
	p.ConfigureParentFlags()

	if errs := p.RunAll(); len(errs) != 0 {
		os.Exit(1)
	}
}
