package main

import (
	"flag"
	"os"
	"path"
)

func main() {
	p := Program{}

	presticDir := ".prestic"

	if homeDir, err := os.UserHomeDir(); err == nil {
		presticDir = path.Join(homeDir, presticDir)
	}

	flag.StringVar(&p.ConfigFile, "config", path.Join(presticDir, "config.yml"), "Config file")
	flag.StringVar(&p.StateFile, "state", path.Join(presticDir, "state.json"), "State file")
	flag.BoolVar(&p.DryRun, "dry-run", false, "Perform a dry run")
	flag.Parse()

	p.LoadConfig()

	if errs := p.RunAll(); len(errs) != 0 {
		os.Exit(1)
	}
}
