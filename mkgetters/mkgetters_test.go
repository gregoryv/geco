package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"testing"
)

func Example_MakeGettersString() {
	src := `package x

type Car struct {
Name string

model string
make string
}`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	out, _ := MakeGetters(file, "Car")
	fmt.Print(string(out))
	// output:
	// func (c *Car) Model() string { return c.model }
	// func (c *Car) Make() string { return c.make }
}

func Test_MakeGettersString(t *testing.T) {
	src := `package x

type Car struct {
Name string

model string
make string
}`

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	// wrong type
	if _, err := MakeGetters(file, "Nosuch"); err == nil {
		t.Error("expect error on missing type")
	}

	if _, err := MakeGetters(file, "Car"); err != nil {
		t.Error(err)
	}
}
