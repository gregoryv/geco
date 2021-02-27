package docs

import "testing"

func Test_generate(t *testing.T) {
	if err := WriteTo("."); err != nil {
		t.Fatal(err)
	}
}
