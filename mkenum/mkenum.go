// Generates get methods for private struct fields
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
	"github.com/gregoryv/geco"
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

	var c geco.MakeEnum
	c.SetFileSet(fset)

	for _, pkg := range packages {
		// should only be one package
		pkgline = fmt.Sprintf("package %s\n", pkg.Name)

		for _, file := range pkg.Files {
			for _, typeName := range typeNames {
				c.SetFile(file)
				c.SetTypeName(typeName)
				if err := c.Run(); err != nil {
					fmt.Println(err)
					continue
				}
				if len(c.Values) == 0 {
					continue
				}
				fmt.Fprint(&buf, "func ")
				fmt.Fprint(&buf, "(")
				fmt.Fprint(&buf, typeName)
				fmt.Fprint(&buf, ") ")
				fmt.Fprint(&buf, "Enum() []any { ")
				buf.WriteString("\n")
				fmt.Fprint(&buf, "	return []any{")
				buf.WriteString("\n	")
				for _, v := range c.Values {
					fmt.Fprint(&buf, v, ",\n")
				}

				fmt.Fprint(&buf, "	}")
				buf.WriteString("\n")
				fmt.Fprint(&buf, "}")
				buf.WriteString("\n")
				c.Values = c.Values[:0]
			}
		}
	}

	// create file
	var header bytes.Buffer
	if appendTo == "" {
		fmt.Fprint(&header, "// GENERATED, DO NOT EDIT!")
		header.WriteString("\n\n")
		header.WriteString(pkgline)
	}
	io.Copy(&header, &buf)

	//fmt.Println(header.String())
	// format it
	nice, err := format.Source(header.Bytes())
	if err != nil {
		log.Fatal(header.String(), "\n", err)
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
