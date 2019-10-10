package main

import (
	"strings"
)

type valueManifestInterface interface {
	open(path string) error
	setValue(path []string, value interface{}) error
	save() error
}

type valuePath struct {
	path  []string
	value interface{}
}

func parseValues(values string) ([]valuePath, error) {
	items := strings.Split(values, ":")
	valueMapList := make([]valuePath, len(items))
	for index, item := range items {
		keyvalue := strings.Split(item, "=")
		valueMapList[index] = valuePath{strings.Split(keyvalue[0], "."), keyvalue[1]}
	}
	return valueMapList, nil
}
