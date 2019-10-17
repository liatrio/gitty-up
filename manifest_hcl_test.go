package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

func TestManifestHclInterface(t *testing.T) {
	manifest := new(manifestHcl)

	assert.Implements(t, (*manifestInterface)(nil), manifest)
}

func TestHclManifestOpen(t *testing.T) {
	// Create HCL file to test
	file, err := ioutil.TempFile("", "test_hcl")
	if err != nil {
		t.Errorf("Test setup failed. Could not create 'test_hcl' file: %s", err)
	}
	defer os.Remove(file.Name())
	file.WriteString("inputs = {\n  one = \"bar\"\n}")
	file.Close()

	hclManifest := &manifestHcl{}
	err = hclManifest.open(file.Name())

	assert.Nil(t, err)
	assert.Equal(t, hclManifest.file, file.Name())
	assert.NotNil(t, hclManifest.astFile)
}

func TestHclManifestSetValue(t *testing.T) {
	value := &ast.LiteralType{Token: token.Token{Type: token.IDENT, Text: "\"foo\""}}
	astFile := &ast.File{
		Node: &ast.ObjectItem{
			Keys: []*ast.ObjectKey{&ast.ObjectKey{Token: token.Token{Type: token.IDENT, Text: "one"}}},
			Val: &ast.ObjectItem{
				Keys: []*ast.ObjectKey{&ast.ObjectKey{Token: token.Token{Type: token.IDENT, Text: "two"}}},
				Val:  value,
			},
		},
	}

	hclManifest := &manifestHcl{"", astFile}
	hclManifest.setValue([]interface{}{"one", "two"}, "bar")

	assert.Equal(t, "\"bar\"", value.Token.Text)
}

func TestHclManifestSave(t *testing.T) {
	value := &ast.LiteralType{Token: token.Token{Type: token.IDENT, Text: "\"foo\""}}
	astFile := &ast.File{
		Node: &ast.ObjectItem{
			Keys: []*ast.ObjectKey{&ast.ObjectKey{Token: token.Token{Type: token.IDENT, Text: "one"}}},
			Val: &ast.ObjectItem{
				Keys: []*ast.ObjectKey{&ast.ObjectKey{Token: token.Token{Type: token.IDENT, Text: "two"}}},
				Val:  value,
			},
		},
	}
	file, err := ioutil.TempFile("", "test_hcl")
	if err != nil {
		t.Errorf("Test setup failed. Could not create 'test_hcl' file: %s", err)
	}
	defer os.Remove(file.Name())

	hclManifest := &manifestHcl{file.Name(), astFile}
	hclManifest.save()

	contents, err := ioutil.ReadFile(file.Name())
	assert.FileExists(t, file.Name())
	assert.NotEmpty(t, contents)
}
