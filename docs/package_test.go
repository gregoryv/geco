package docs

import "testing"

func Test_generate(t *testing.T) {
	if err := Generate("."); err != nil {
		t.Fatal(err)
	}
}
