package main

import (
	"strings"
)

type valueManifestInterface interface {
	open(path string) error
	setValue(path []interface{}, value interface{}) error
	save() error
}

type valuePath struct {
	path  []interface{}
	value interface{}
}

func parseValues(values string) ([]valuePath, error) {
	items := strings.Split(values, ":")
	valueMapList := make([]valuePath, len(items))
	for index, item := range items {
		keyvalue := strings.Split(item, "=")
		crumbs := strings.Split(keyvalue[0], ".")
		path := make([]interface{}, len(crumbs))
		for k, v := range crumbs {
			path[k] = v
		}
		valueMapList[index] = valuePath{path, keyvalue[1]}
	}
	return valueMapList, nil
}
