package main

import (
	"testing"
	"io/ioutil"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestManifestJSONInterface(t *testing.T) {
	manifest := new(manifestJSON)

	assert.Implements(t, (*manifestInterface)(nil), manifest)
}

func TestManifestJSONOpen(t *testing.T) {

	manifest := manifestJSON{}
	err := manifest.open("./sample/sample_a.json")
	assert.NoError(t, err)

	assert.Equal(t, "v0.1.1", manifest.data.(map[string]interface{})["one"].(map[string]interface{})["magenta"])
}

func TestManifestJSONSetValue(t *testing.T) {
	manifest := manifestJSON{}
	manifest.data = map[string]interface{}{"one": map[string]interface{}{"cyan": "v0.0.1", "magenta": "v0.1.1"}, "two": map[string]interface{}{"yellow": "v0.2.1", "black": "v0.3.1"}}

	err := manifest.setValue([]interface{}{"one", "magenta"}, "v0.0.42")
	assert.NoError(t, err)

	assert.Equal(t, "v0.0.42", manifest.data.(map[string]interface{})["one"].(map[string]interface{})["magenta"])
}

func TestManifestJSONSave(t *testing.T) {
	file, err := ioutil.TempFile("", "test_json_save")
	assert.NoError(t, err)
	file.Close()

	manifest := manifestJSON{}
	manifest.file = file.Name()
	manifest.data = map[string]interface{}{"one": map[string]interface{}{"cyan": "v0.0.1", "magenta": "v0.1.1"}, "two": map[string]interface{}{"yellow": "v0.2.1", "black": "v0.3.1"}}

	err = manifest.save();
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile("./sample/sample_a.json")
	assert.NoError(t, err)

	actual, err := ioutil.ReadFile(file.Name())
	assert.NoError(t, err)

	// assert.Equal(t, expected, actual)
	assert.Less(t, 0, len(strings.TrimSpace(string(actual))))
	assert.Equal(t, len(strings.TrimSpace(string(expected))), len(strings.TrimSpace(string(actual))))
}