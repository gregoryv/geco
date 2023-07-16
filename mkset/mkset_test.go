package main

import (
	"io"
	"os/exec"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_mkset(t *testing.T) {
	out, err := exec.Command("go", "generate", ".").CombinedOutput()
	if err != nil {
		t.Log(string(out))
		t.Fatal(err)
	}
	golden.Assert(t, string(out))
}

//go:generate mkset -t Car,Boat
type Car struct {
	Name string

	model  string
	make   int
	output io.Writer
}

type Boat struct {
	color int
}
