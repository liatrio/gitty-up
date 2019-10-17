package main

import (
	"testing"
	"io/ioutil"
	"os"

	"github.com/stretchr/testify/assert"

)

func TestManifestJSONInterface(t *testing.T) {
	manifest := new(manifestJSON)

	assert.Implements(t, (*manifestInterface)(nil), manifest)
}

func TestManifestJSONOpen(t *testing.T) {
	file, err := ioutil.TempFile("", "test_json_open")
	if err != nil {
		t.Errorf("Test setup failed. Could not create test_json_open file: %s", err)
	}
	defer os.Remove(file.Name())
	file.WriteString("{\"inputs\": {\"one\": \"bar\"}}")
	file.Close()

	manifest := &manifestJSON{}
	err = manifest.open(file.Name())

	assert.Nil(t, err)
	assert.Equal(t, manifest.file, file.Name())
	assert.IsType(t, make(map[string]interface{}), manifest.data)
	assert.Equal(t, manifest.data.(map[string]interface{})["inputs"].(map[string]interface{})["one"], "bar")
}

func TestManifestJSONSetValue(t *testing.T) {
	jsonObj := &map[string]interface{}{"one": map[string]interface{}{"two": "foo"}}
	// jsonObj := &map[string]map[string]string{}
	path := []interface{}{"one", "two"}

	// (*jsonObj)["one"] = map[string]string{"two": "foo"}

	manifest := &manifestJSON{"", *jsonObj}
	err := manifest.setValue(path, "bar")

	value := (*jsonObj)["one"].(map[string]interface{})["two"]
	assert.Nil(t, err)
	assert.Equal(t, "bar", value)
}

func TestManifestJSONSave(t *testing.T) {
	file, err := ioutil.TempFile("", "test_json_save")
	if err != nil {
		t.Errorf("Test setup failed. Could not create test_json_save")
	}
	defer os.Remove(file.Name())
	file.Close()

	jsonObj := &map[string]interface{}{"one": map[string]interface{}{"two": "foo"}}
	
	manifest := &manifestJSON{file.Name(), jsonObj}

	err = manifest.save();
	assert.Nil(t, err)

	contents, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Error("Test setup failed. Could not read output file")
	}
	assert.FileExists(t, file.Name())
	assert.Equal(t, "{\n  \"one\": {\n    \"two\": \"foo\"\n  }\n}\n", string(contents))
}