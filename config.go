package main

import "gopkg.in/yaml.v2"

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
