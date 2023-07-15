package main

import (
	"bytes"
	"fmt"
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
				data, err := MakeGetters(file, typeName)
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
	var result bytes.Buffer
	fmt.Fprint(&result, "// GENERATED!, DO NOT EDIT!")
	result.WriteString("\n\n")
	result.WriteString(pkgline)
	io.Copy(&result, &buf)

	// format it
	nice, err := format.Source(result.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// write result
	switch {
	case saveTo == "":
		os.Stdout.Write(nice)

	default:
		if err := os.WriteFile(saveTo, nice, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
