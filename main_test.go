package main

import (
  _ "os"
  "testing"
  _ "reflect"
  "github.com/stretchr/testify/assert"
  _ "github.com/hashicorp/hcl"
  _ "github.com/hashicorp/hcl/hcl/printer"
)

// func TestDecodeHcl(t *testing.T) {
//   hclString := "inputs = {\n  one = \"bar\"\n}"
//   f, err := decodeHcl(hclString)

//   assert.Equal(t, reflect.TypeOf(*f).Name(), "File")

//   assert.Nil(t, err)
// }

// func TestSetValueInAst(t *testing.T) {
//   hclString := "inputs = {\n  one = \"bar\"\n}"

//   ast, _ := hcl.Parse(hclString)

//   //First test runs correct input into setValueInAst to change 'bar' to 'foo'
//   err := setValueInAst(*ast, []string{"inputs", "one"}, "foo")

//   file, err := os.Create("./testsetvalue.hcl")
//   err = printer.Fprint(file, ast.Node)

//   fileInfo, err := file.Stat()

//   data := make([]byte, fileInfo.Size())
//   file.Seek(0,0)
//   _, err = file.Read(data)
//   hclStringResult := "inputs = {\n  one = \"foo\"\n}"

//   assert.Equal(t, string(data), hclStringResult)

//   assert.Nil(t, err)

//   os.Remove("./testsetvalue.hcl")

//   //Second test tries to access a value that doesn't exist 'two'
//   err = setValueInAst(*ast, []string{"inputs", "missing", "two"}, "foo")

//   assert.Equal(t, err.Error(), "Did not match value ([missing two] -> foo)")
// }

func TestParseValues(t *testing.T) {
  values := "input.one=foo"
  valuePath, err := parseValues(values)

  assert.Nil(t, err)

  assert.Equal(t, valuePath[0].path[0].(string), "input")
  assert.Equal(t, valuePath[0].path[1].(string), "one")
  assert.Equal(t, valuePath[0].value.(string), "foo")
}

