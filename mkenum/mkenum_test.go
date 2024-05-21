package main

import (
	"os/exec"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_mkenum(t *testing.T) {
	out, _ := exec.Command("go", "generate", ".").Output()
	golden.Assert(t, string(out))
}
