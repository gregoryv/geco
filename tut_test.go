package geco

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGenerateShould(t *testing.T) {
	gc := NewTypeUnderTest("x", "CarUnderTest", "Car", "", []byte(src))
	var buf bytes.Buffer
	gc.Generate(&buf)

	got := buf.String()
	if strings.Contains(got, "shouldModel") {
		t.Errorf("contains non failure possible:\n%s", got)
	}
	exp := []string{
		"package x",
		"shouldWorks",
		"shouldServe(w *http.Request)",
		"shouldYWorks(a, b int, v string)",
		") UnderTest(t *testing.T) *CarUnderTest",

		"err := me.Serve(w)",
		"err := me.Works()",
		"err := me.XWorks(a, b)",
		"err := me.YWorks(a, b, v)",
		"v0, err := me.KindaWorks(v)",
		"if err != nil",
		"me.T.Helper()",
		"me.T.Error(err)",
		"*testing.T",
		"return v0",
	}
	for _, exp := range exp {
		if !strings.Contains(got, exp) {
			t.Fatalf("missing %q:\n%s", exp, got)
		}
	}
}

// ----------------------------------------

const src = `package example
// todo: implement more
func main() {}
type Car struct{}
func (me *Car) Model() string { return "norAx" }
func (me *Car) Age() int      { return 0 }
func (me *Car) Serve(w *http.Request) error { return nil }
func (me Car) Works() error { return nil }
func (me *Car) XWorks(a,b int) error { return nil }
func (me *Car) YWorks(a,b int, v string) error { return nil }
func (me *Car) KindaWorks(v int) (int, error) { return 0, nil }
func (me *Car) JustDoIt() { }
func Model() error { return nil}
func (me *Boat) FixMe() error { return nil }
`

// ----------------------------------------

func mustParse(t *testing.T, src string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	mode := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, "", src, mode)
	if err != nil {
		t.Fatal(err)
	}
	return fset, file
}

func toString(fset *token.FileSet, file *ast.File) string {
	var buf bytes.Buffer
	format.Node(&buf, fset, file)
	return buf.String()
}
