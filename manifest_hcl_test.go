package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

func TestManifestHclInterface(t *testing.T) {
	manifest := new(manifestHcl)

	assert.Implements(t, (*manifestInterface)(nil), manifest)
}

func TestHclManifestOpen(t *testing.T) {
	manifest := &manifestHcl{}
	err := manifest.open("./sample/sample_a.hcl")
	assert.NoError(t, err)

	assert.NotNil(t, manifest.astFile)
	assert.Equal(t, "\"v0.1.1\"", manifest.astFile.Node.(*ast.ObjectList).Items[0].Val.(*ast.ObjectType).List.Items[1].Val.(*ast.LiteralType).Token.Text)
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
	assert.NoError(t, err)

	hclManifest := &manifestHcl{file.Name(), astFile}
	err = hclManifest.save()
	assert.NoError(t, err)

	contents, err := ioutil.ReadFile(file.Name())
	assert.NoError(t, err)
	assert.FileExists(t, file.Name())
	assert.NotEmpty(t, contents)
}
