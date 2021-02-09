package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test_gentut_to_stdout(t *testing.T) {
	os.Chdir("./testdata")
	out, err := exec.Command("go", "run", "..",
		"-in", "car.go", "-p", "testdata", "-t", "Car",
	).CombinedOutput()
	if err != nil {
		t.Fatal(err, string(out))
	}
	got := string(out)

	// Default should write to stdout
	if !strings.Contains(got, "package testdata") {
		t.Error(got)
	}
}

func Test_gentut_to_file(t *testing.T) {
	os.Chdir("./testdata")
	out, err := exec.Command("go", "run", "..",
		"-in", "car.go", "-p", "testdata", "-t", "Car", "-w",
	).CombinedOutput()
	if err != nil {
		t.Fatal(err, string(out))
	}
	got := string(out)

	// Default should not write to stdout
	if strings.Contains(got, "package testdata") {
		t.Error(got)
	}
	// File must be created
	if err := os.Remove("carut_test.go"); err != nil {
		t.Fatal(err)
	}
}

func Test_outFilename(t *testing.T) {
	cases := map[string]string{
		"name.go":  "nameut_test.go",
		"./x.go":   "xut_test.go",
		"dir/y.go": "dir/yut_test.go",
	}
	for name, exp := range cases {
		t.Run(name, func(t *testing.T) {
			got := outFilename(name)
			if got != exp {
				t.Errorf("got %q, expected %q", got, exp)
			}
		})
	}

}
