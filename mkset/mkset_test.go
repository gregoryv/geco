package main

import (
	"os"
	"os/exec"
)

func Example_mkset() {
	// see testdata/example.go for the go:generate statement
	out, _ := exec.Command("go", "generate", ".").Output()
	os.Stdout.Write(out)
	// output:
	// // GENERATED, DO NOT EDIT!
	//
	// package main
	//
	// func (c *Car) SetModel(v string) { c.model = v }
	// func (c *Car) SetMake(v int)     { c.make = v }
	//
	// func (b *Boat) SetColor(v int) { b.color = v }
}

//go:generate mkset -t Car,Boat
type Car struct {
	Name string

	model string
	make  int
}

type Boat struct {
	color int
}
