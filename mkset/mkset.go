// Generates get methods for private struct fields
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gregoryv/cmdline"
)

func main() {
	log.SetFlags(0)
	var (
		cli = cmdline.NewBasicParser()

		types    = cli.Option("-t, --types", "CSV list of types").String("")
		saveTo   = cli.Option("-w, --write-file").String("")
		appendTo = cli.Option("-a, --append-file").String("")
		filename = cli.NamedArg("FILE").String(".")
	)
	cli.Parse()

	// validate options
	if filename == "" || types == "" {
		cli.Usage().WriteTo(os.Stderr)
		os.Exit(1)
	}

	var (
		buf       bytes.Buffer
		pkgline   string
		typeNames = strings.Split(types, ",")
	)

	// prepare filter
	path := filename
	var name string
	if strings.Contains(filename, string(filepath.Separator)) {
		path, name = filepath.Split(filename)
	}
	filter := func(f fs.FileInfo) bool {
		if strings.HasSuffix(name, ".go") {
			return f.Name() == name
		}
		return true
	}

	// parse source
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, path, filter, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range packages {
		// should only be one package
		pkgline = fmt.Sprintf("package %s\n", pkg.Name)

		for _, file := range pkg.Files {
			for _, typeName := range typeNames {
				data, err := MakeSetters(fset, file, typeName)
				if err != nil {
					continue
				}
				buf.Write(data)
				buf.WriteString("\n")
			}
		}
	}

	if buf.Len() == 0 {
		log.Fatalf("error: no private fields found for %s in %s", types, path)
	}

	// create file
	var header bytes.Buffer

	if appendTo == "" {
		fmt.Fprint(&header, "// GENERATED, DO NOT EDIT!")
		header.WriteString("\n\n")
		header.WriteString(pkgline)
	}
	io.Copy(&header, &buf)

	// format it
	nice, err := format.Source(header.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// write result
	switch {
	case saveTo != "":
		if err := os.WriteFile(saveTo, nice, 0644); err != nil {
			log.Fatal(err)
		}

	case appendTo != "":
		f, err := os.OpenFile(appendTo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		f.Write(nice)
		f.Close()

	default:
		os.Stdout.Write(nice)
	}
}

// MakeSetters returns generated set methods for the given type.
func MakeSetters(fset *token.FileSet, file *ast.File, typeName string) ([]byte, error) {
	// find fields
	var fields []*ast.Field
	var name string
	ast.Inspect(file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			name = n.Name // save name

		case *ast.StructType:
			if name == typeName {
				fields = n.Fields.List
			}
		}
		return true
	})

	if len(fields) == 0 {
		return nil, fmt.Errorf("type %s: not found", typeName)
	}
	var buf bytes.Buffer
	for _, field := range fields {
		name := field.Names[0].Name
		if !isPrivate(name) {
			continue
		}
		rcv := receiver(typeName)

		fmt.Fprint(&buf, "func ")
		fmt.Fprint(&buf, "(")
		fmt.Fprint(&buf, rcv)
		fmt.Fprint(&buf, " ")
		fmt.Fprint(&buf, "*")
		fmt.Fprint(&buf, typeName)
		fmt.Fprint(&buf, ") ")
		fmt.Fprint(&buf, "Set")
		fmt.Fprint(&buf, makePublic(name))
		fmt.Fprint(&buf, "(v ")
		format.Node(&buf, fset, field.Type)
		fmt.Fprint(&buf, ") ")
		fmt.Fprint(&buf, "{ ")
		fmt.Fprint(&buf, rcv)
		fmt.Fprint(&buf, ".")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, " = v ")
		fmt.Fprint(&buf, "}")
		fmt.Fprintln(&buf)
	}

	return buf.Bytes(), nil
}
