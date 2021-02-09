package geco

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strings"

	"github.com/gregoryv/nexus"
)

func NewGeneratorFrom(pkg, rec, typ string, src interface{}) *Generator {
	return NewGenerator(pkg, rec, typ, "", src)
}

func NewGenerator(pkg, rec, typ, filename string, src interface{}) *Generator {
	return &Generator{
		Package:  pkg,
		Receiver: rec,
		Type:     typ,
		filename: filename,
		src:      src,
	}
}

type Generator struct {
	Package  string
	Receiver string
	Type     string

	filename string
	src      interface{}

	prefix     string
	errHandler string // e.g. Error or Fatal
	p          *nexus.Printer
	fs         *token.FileSet
}

// Generate parses the given filename or src and generates a type with
// should/must test helper methods for any method that returns an
// error.  Methods prefixed with should, use t.Error as error handler
// and those prefixed with must use t.Fatal.
func (me *Generator) Generate(w io.Writer) error {
	fset := token.NewFileSet()
	mode := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, me.filename, me.src, mode)
	if err != nil {
		return err
	}
	p, _ := nexus.NewPrinter(w)
	p.Println("package", me.Package)
	p.Println()
	p.Println(`import "testing"`)
	p.Println()

	// It's possible to reuse the existing type as the receiver.
	if me.Receiver != me.Type {
		p.Printf("type %s struct {\n", me.Receiver)
		p.Println("\t*testing.T")
		p.Printf("\t*%s\n", me.Type)
		p.Println("}")
		p.Println()
	}

	// Add geco(testing.T) constructor on the existing type, conflict
	// could occur.
	p.Printf("func (me *%s) tu(t *testing.T) *%s {\n", me.Type, me.Receiver)
	p.Printf("\treturn &%s{T:t, %s: me}\n", me.Receiver, me.Type)
	p.Println("}")
	p.Println()

	me.fs = fset
	me.p = p
	ast.Inspect(file, me.visit)
	return nil
}

func (me *Generator) visit(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.FuncDecl:
		if n.Type.Results == nil { // skip
			return true
		}
		if !returnsError(n) { // skip
			return true
		}
		if n.Recv == nil { // skip
			return true
		}
		rname := sprintType(me.fs, n.Recv.List[0].Type)
		if rname != me.Type && rname != "*"+me.Type {
			return true
		}

		me.printFunc(n, "should", "Error")
		me.printFunc(n, "must", "Fatal")
	}
	return true
}

func (me *Generator) printFunc(n *ast.FuncDecl, prefix, errhand string) {
	p := me.p
	// method receiver and name
	p.Printf("func (me *%s) %s%s(", me.Receiver, prefix, n.Name)

	// method parameters
	params := make([]string, 0)
	for _, field := range n.Type.Params.List {
		args := make([]string, 0)
		for _, n := range field.Names {
			args = append(args, fmt.Sprint(n))
		}
		namedParam := fmt.Sprintf(
			"%s %s", csv(args), sprintType(me.fs, field.Type),
		)
		params = append(params, namedParam)
	}
	p.Print(csv(params))
	p.Print(")")

	// signature returns values, except error
	var buf bytes.Buffer
	sign, _ := nexus.NewPrinter(&buf)
	for _, field := range n.Type.Results.List {
		v := sprintType(me.fs, field.Type)
		if v != "error" {
			sign.Print(" ", v)
		}
	}
	p.Print(buf.String())

	// body of method
	p.Println(" {")
	rv := strings.Join(retVals(n), ", ")
	args := make([]string, 0)
	for _, field := range n.Type.Params.List {
		for _, n := range field.Names {
			args = append(args, fmt.Sprint(n))
		}
	}
	p.Printf("\t%s := me.%s(%s)\n", rv, n.Name, csv(args))
	p.Println("\tif err != nil {")
	p.Println("\t\tme.T.Helper()")
	p.Printf("\t\tme.T.%s(err)\n", errhand)
	p.Println("\t}")
	l := len(n.Type.Results.List)
	if l > 1 {
		p.Print("\treturn ")
		retvals := make([]string, l-1)
		for i, _ := range n.Type.Results.List[:l-1] {
			retvals[i] = fmt.Sprintf("v%d", i)
		}
		p.Println(strings.Join(retvals, ", "))
	}
	p.Println("}")
	p.Println()
}

func sprintType(fset *token.FileSet, n ast.Expr) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, n)
	return buf.String()
}

func csv(v []string) string {
	return strings.Join(v, ", ")
}

func returnsError(n *ast.FuncDecl) bool {
	for _, field := range n.Type.Results.List {
		if fmt.Sprint(field.Type) == "error" {
			return true
		}
	}
	return false
}

func retVals(n *ast.FuncDecl) []string {
	vals := make([]string, 0)
	for i, field := range n.Type.Results.List {
		if fmt.Sprint(field.Type) == "error" {
			vals = append(vals, "err")
			continue
		}
		vals = append(vals, fmt.Sprintf("v%d", i))
	}
	return vals
}

type modFunc func(ast.Node) bool
