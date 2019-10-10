package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/printer"
)

type manifestHcl struct {
	file    string
	astFile *ast.File
}

func (h *manifestHcl) open(file string) (err error) {
	h.file = file

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	fmt.Println("Decoding HCL data")
	h.astFile, err = hcl.Parse(string(contents))
	if err != nil {
		return
	}

	return
}

func (h *manifestHcl) setValue(path []string, value interface{}) error {
	fmt.Printf("Changing value of %v\n", path)
	matched := false
	ast.Walk(h.astFile.Node, func(n ast.Node) (ast.Node, bool) {
		if item, ok := n.(*ast.ObjectItem); ok {
			for _, key := range item.Keys {
				if key.Token.Type.IsIdentifier() && key.Token.Text == path[0] {
					if len(path) == 1 {
						if val, ok := item.Val.(*ast.LiteralType); ok {
							fmt.Printf("Changed value %s -> \"%s\"\n", item.Val.(*ast.LiteralType).Token.Text, value)
							val.Token.Text = fmt.Sprintf("\"%s\"", value)
							matched = true
						} else {
							fmt.Printf("Warning: Cannot change value of type %T\n", item.Val)
						}
						return n, false // we matched the end of the path
					}
					path = path[1:]
					return n, true // we matched part of the path
				}
				return n, false // this branch does not match our path
			}
		}
		return n, true // traverse all non item nodes
	})
	if matched == false {
		return fmt.Errorf("Did not match value (%v -> %s)", path, value)
	}
	return nil

}

func (h *manifestHcl) save() (err error) {
	file, err := os.Create(h.file)
	if err != nil {
		return
	}

	defer file.Close()
	err = printer.Fprint(file, h.astFile.Node)

	return
}