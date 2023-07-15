package goref

import (
	"fmt"
	"testing"
)

func Example_MakeGettersString() {
	src := `package x

type Car struct {
Name string

model string
make string
}`
	out, _ := MakeGettersString(src, "Car")
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

	// wrong type
	if _, err := MakeGettersString(src, "Nosuch"); err == nil {
		t.Error("expect error on missing type")
	}

	if _, err := MakeGettersString(src, "Car"); err != nil {
		t.Error(err)
	}
}
