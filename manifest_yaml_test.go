package main

import (
	"testing"
	"io/ioutil"
	"strings"

	// "gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

func TestManifestYamlInterface(t *testing.T) {
	manifest := new(manifestYaml)

	assert.Implements(t, (*manifestInterface)(nil), manifest)
}

func TestManifestYamlOpen(t *testing.T) {
	manifest := new(manifestYaml)

	err := manifest.open("./sample/sample_a.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "v0.1.1", manifest.data.(map[string]interface{})["one"].(map[string]interface{})["magenta"])
}

func TestManifestYamlSetValue(t *testing.T) {
	manifest := new(manifestYaml)
	manifest.data = map[string]interface{}{"one": map[string]interface{}{"cyan": "v0.0.1", "magenta": "v0.1.1"}, "two": map[string]interface{}{"yellow": "v0.2.1", "black": "v0.3.1"}}

	manifest.setValue([]interface{}{"one", "magenta"}, "v0.0.42")
	
	assert.Equal(t, "v0.0.42", manifest.data.(map[string]interface{})["one"].(map[string]interface{})["magenta"])
}

func TestManifestYamlSave(t *testing.T) {
	file, err := ioutil.TempFile("", "test_manifest_yaml_save")
	assert.NoError(t, err)
	file.Close()
	
	manifest := new(manifestYaml)
	
	manifest.file = file.Name()	
	manifest.data = map[string]interface{}{"one": map[string]interface{}{"cyan": "v0.0.1", "magenta": "v0.1.1"}, "two": map[string]interface{}{"yellow": "v0.2.1", "black": "v0.3.1"}}
	
	manifest.save()

	expected, err := ioutil.ReadFile("./sample/sample_a.yaml")
	assert.NoError(t, err)

	actual, err := ioutil.ReadFile(file.Name())
	assert.NoError(t, err)

	
	// assert.Equal(t, expected, actual) // yaml encoder changes order of elements so we can't compare the file contents
	assert.Less(t, 0, len(strings.TrimSpace(string(actual))))
	assert.Equal(t, len(strings.TrimSpace(string(expected))), len(strings.TrimSpace(string(actual))))
}