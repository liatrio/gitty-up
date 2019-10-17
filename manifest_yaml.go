package main

import (
	"fmt"
	"os"
	"io"

	"gopkg.in/yaml.v3"
)

type manifestYaml struct {
	file string
	data interface{}
}

func (m *manifestYaml) open(file string) (err error) {
	m.file = file

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Failed opening YAML: %s", err)
	}
	decoder := yaml.NewDecoder(f)
	for {
		err = decoder.Decode(&m.data)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("Failed to parse YAML: %s", err)
		}
	}

}

func (m *manifestYaml) setValue(path []interface{}, value interface{}) (err error) {
	fmt.Printf("Setting value %v: %v\n", path, value)
	return walk(&m.data, path, value)
}

func (m *manifestYaml) save() (err error) {
	f, err := os.Create(m.file)
	if err != nil {
		return
	}
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	err = encoder.Encode(m.data)
	if err != nil {
		return
	}
	
	return
}