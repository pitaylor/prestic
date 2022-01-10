package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

// GetCommand gets a Command by name.
func (c *Config) GetCommand(name string) (*Command, error) {
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return &cmd, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("command not found: %v", name))
}

// UnmarshalYAML unmarshalls a config ignoring keys with "x-" prefix that can be used for yaml anchors.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var config struct {
		Commands CommandList            `yaml:",omitempty"`
		Rest     map[string]interface{} `yaml:",inline"`
	}

	if err := unmarshal(&config); err != nil {
		return err
	}

	for key := range config.Rest {
		if !strings.HasPrefix(key, "x-") {
			return errors.New(fmt.Sprintf("unknown config property \"%v\"", key))
		}
	}

	c.Commands = config.Commands

	return nil
}

// UnmarshalYAML unmarshalls a map of Flags as an array, preserving the ordering specified by the YAML.
func (l *FlagList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var mapSlice yaml.MapSlice

	if err := unmarshal(&mapSlice); err != nil {
		return err
	}

	for _, mapItem := range mapSlice {
		name := fmt.Sprintf("%v", mapItem.Key)
		*l = append(*l, Flag{Name: name, Value: mapItem.Value})
	}

	return nil
}

// UnmarshalYAML unmarshalls a map of Commands as an array, preserving the ordering specified by the YAML.
func (l *CommandList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cmdMap map[string]*Command

	if err := unmarshal(&cmdMap); err != nil {
		return err
	}

	var mapSlice yaml.MapSlice

	if err := unmarshal(&mapSlice); err != nil {
		return err
	}

	for _, mapItem := range mapSlice {
		name := mapItem.Key.(string)
		cmdMap[name].Name = name
		*l = append(*l, *cmdMap[name])
	}

	return nil
}
