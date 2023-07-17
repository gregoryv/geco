package main

import (
	"io"
	"os/exec"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_mkfill(t *testing.T) {
	out, _ := exec.Command("go", "generate", ".").Output()
	golden.Assert(t, string(out))
}

//go:generate mkfill -t Car,Boat
type Car struct {
	Name string

	model  string
	make   int
	output io.Writer
}

type Boat struct {
	color int
	model string // same
}
