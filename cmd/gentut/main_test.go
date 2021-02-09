package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func Test_run_gentest(t *testing.T) {
	os.Chdir("./testdata")
	defer os.RemoveAll("carsut_test.go")
	out, err := exec.Command(
		"go", "run", "..",
		"-in", "cars.go", "-p", "testdata", "-t", "Car",
	).CombinedOutput()
	if err != nil {
		t.Fatal(err, string(out))
	}

	body, err := ioutil.ReadFile("carsut_test.go")
	if err != nil {
		t.Fatal(err)
	}
	t.Error(string(body))

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
