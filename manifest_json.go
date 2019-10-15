package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"encoding/json"
)

type manifestJSON struct {
	file string
	jsonObj interface{}
}

func (m *manifestJSON) open(file string) (err error) {
	m.file = file

	f, err := os.Open(m.file)
	if err != nil {
		return fmt.Errorf("Error opening JSON file: %s", err)
	}

	decoder := json.NewDecoder(f)

	for {
		err := decoder.Decode(&m.jsonObj)
		if err == io.EOF {
			break;
		}
		if err != nil {
			return fmt.Errorf("Error reading JSON file: %s", err)
		}
	}

	return
}

func (m *manifestJSON) setValue(path []interface{}, value interface{}) (error) {
	fmt.Printf("Setting value %v: %v\n", path, value)
	return walk(&m.jsonObj, path, value)
}

func (m *manifestJSON) save() (err error) {
	f, err := os.Create(m.file)
	if err != nil {
		return fmt.Errorf("Error opening JSON file: %s", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	encoder.SetIndent("", "  ")
	err = encoder.Encode(m.jsonObj);
	if err != nil {
		return fmt.Errorf("Error encoding JSON to file: %s", err)
	}

	return
}

func walk(pos *interface{}, path[]interface{}, value interface{}) error {
	switch reflect.TypeOf(*pos).Kind() {
	case reflect.Map:
		for k, v := range (*pos).(map[string]interface{}) {
			if k == path[0] {
				if (len(path) == 1) {
					(*pos).(map[string]interface{})[k] = value
					fmt.Printf("Changed value %s -> %s\n", v, value)
					return nil
				}
				result := walk(&v, path[1:], value)
				return result
			}
		}
		return fmt.Errorf("Set value failed. Could not find node '%s' in path", path[0])
	default:
		return fmt.Errorf("Set value failed. Unhandled type '%t' for node '%s'", *pos, path[0])
	}
}