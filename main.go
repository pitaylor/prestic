package main

import "flag"

func main() {
	//snapshot c098c88e saved

	p := Program{}

	flag.StringVar(&p.ConfigFile, "config", "prestic.yml", "Config file")
	flag.StringVar(&p.StateFile, "state", ".prestic.json", "State file")
	flag.BoolVar(&p.DryRun, "dry-run", false, "Perform a dry run")
	flag.Parse()

	p.LoadConfig()
	p.LoadState()
	p.RunAll()
}
