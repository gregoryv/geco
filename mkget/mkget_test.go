package main

import (
	"io"
	"os/exec"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_mkget(t *testing.T) {
	out, _ := exec.Command("go", "generate", ".").Output()
	golden.Assert(t, string(out))
}

//go:generate mkget -t Car
type Car struct {
	Name string

	model  string
	make   int
	output io.Writer
}

type Boat struct {
	color int
}
