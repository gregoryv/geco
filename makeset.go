package goref

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
)

type MakeSet struct {
	fileSet  *token.FileSet
	file     *ast.File
	typeName string
	output   io.Writer
}

func (m *MakeSet) Run() error {
	// find fields
	var fields []*ast.Field
	var name string
	ast.Inspect(m.file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			name = n.Name // save name

		case *ast.StructType:
			if name == m.typeName {
				fields = n.Fields.List
			}
		}
		return true
	})

	if len(fields) == 0 {
		return fmt.Errorf("type %s: not found", m.typeName)
	}
	var buf bytes.Buffer
	for _, field := range fields {
		name := field.Names[0].Name
		if !isPrivate(name) {
			continue
		}
		rcv := receiver(m.typeName)

		fmt.Fprint(&buf, "func ")
		fmt.Fprint(&buf, "(")
		fmt.Fprint(&buf, rcv)
		fmt.Fprint(&buf, " ")
		fmt.Fprint(&buf, "*")
		fmt.Fprint(&buf, m.typeName)
		fmt.Fprint(&buf, ") ")
		fmt.Fprint(&buf, "Set")
		fmt.Fprint(&buf, makePublic(name))
		fmt.Fprint(&buf, "(v ")
		format.Node(&buf, m.fileSet, field.Type)
		fmt.Fprint(&buf, ") ")
		fmt.Fprint(&buf, "{ ")
		fmt.Fprint(&buf, rcv)
		fmt.Fprint(&buf, ".")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, " = v ")
		fmt.Fprint(&buf, "}")
		fmt.Fprintln(&buf)
	}
	_, err := io.Copy(m.output, &buf)
	return err
}
