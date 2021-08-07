package main

import (
	"flag"
)

func main() {
	p := Program{}

	flag.StringVar(&p.ConfigFile, "config", "prestic.yml", "Config file")
	flag.BoolVar(&p.DryRun, "dry-run", false, "Perform a dry run")
	flag.Parse()

	p.Load()
	p.Run()
}
