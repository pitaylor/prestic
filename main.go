package main

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var log = logrus.New()

type RunCmd struct {
	Commands []string `arg:"" name:"command" optional:"" help:"Restic commands to run. Runs all if unspecified."`
}

func (r *RunCmd) Run(p *Program) error {
	// Run all commands by default
	commandList := p.Config.Commands

	if len(r.Commands) != 0 {
		commandList = CommandList{}
		for _, name := range r.Commands {
			if command, err := p.Config.GetCommand(name); err == nil {
				commandList = append(commandList, *command)
			} else {
				return err
			}
		}
	}

	var failed []string
	for _, command := range commandList {
		if p.Run(command) != nil {
			failed = append(failed, command.Name)
		}
	}

	if len(failed) != 0 {
		return errors.New(fmt.Sprintf("one or more commands failed: %v", strings.Join(failed, ", ")))
	}

	return nil
}

type ListCmd struct {
}

func (l *ListCmd) Run(p *Program) error {
	fmt.Println("Restic commands:")
	for _, cmd := range p.Config.Commands {
		fmt.Printf("  %-15s\t%v %v\n", cmd.Name, cmd.Command, strings.Join(cmd.Args, " "))
	}
	return nil
}

type CLI struct {
	DryRun     bool   `help:"Enable dry-run mode."`
	ConfigFile string `help:"Configuration file." short:"c" type:"path" default:"${configFile}"`
	StateFile  string `help:"State file." short:"s" type:"path" default:"${stateFile}"`
	Log        struct {
		Level string `help:"Log level: ${enum}." enum:"debug,info,warn,error" default:"info"`
	} `embed:"" prefix:"log-"`

	Run  RunCmd  `cmd:"" help:"Run restic commands."`
	List ListCmd `cmd:"" help:"List restic commands."`
}

// ConfigPaths determines default config and state file paths based on searching for a config file in user
// directories first, then system directories. It loosely follows XDG conventions.
func ConfigPaths() []string {
	var candidates [][]string

	if homeDir, err := os.UserHomeDir(); err == nil {
		candidates = append(
			candidates,
			[]string{homeDir + "/.prestic/config.yml", homeDir + "/.local/share/prestic/state.json"},
			[]string{homeDir + "/.config/prestic/config.yml", homeDir + "/.local/share/prestic/config.yml"},
		)
	}

	candidates = append(
		candidates,
		[]string{"/etc/prestic/config.yml", "/var/prestic/state.json"},
		[]string{"/usr/local/etc/prestic/config.yml", "/usr/local/var/prestic/state.json"},
	)

	for i := range candidates {
		if src, err := os.Stat(candidates[i][0]); err == nil && src.Mode().IsRegular() {
			return candidates[i]
		}
	}

	return candidates[0]
}

func main() {
	var cli CLI

	defaultPaths := ConfigPaths()

	ctx := kong.Parse(&cli,
		kong.Vars{
			"configFile": defaultPaths[0],
			"stateFile":  defaultPaths[1],
		})

	p, err := NewProgram(&cli)

	if err == nil {
		err = ctx.Run(p)
	}

	ctx.FatalIfErrorf(err)
}
