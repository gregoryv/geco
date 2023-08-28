package geco

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
)

type MakeNoZero struct {
	fileSet  *token.FileSet
	file     *ast.File
	typeName string
	output   io.Writer
}

func (m *MakeNoZero) Run() error {
	// find fields
	var fields []*ast.Field
	var name string
	ast.Inspect(m.file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			name = n.Name // save name

		case *ast.StructType:
			if name == m.typeName {
				for _, field := range n.Fields.List {
					name := field.Names[0].Name
					if !isPrivate(name) {
						continue
					}
					fields = append(fields, field)
				}
			}
		}
		return true
	})

	if len(fields) == 0 {
		return fmt.Errorf("type %s: not found", m.typeName)
	}
	rcv := receiver(m.typeName)

	var buf bytes.Buffer
	fmt.Fprint(&buf, "// NoZero returns error if any private field is zero")
	fmt.Fprintln(&buf)
	fmt.Fprint(&buf, "func ")
	fmt.Fprint(&buf, "(")
	fmt.Fprint(&buf, rcv)
	fmt.Fprint(&buf, " ")
	fmt.Fprint(&buf, "*")
	fmt.Fprint(&buf, m.typeName)
	fmt.Fprint(&buf, ") ")
	fmt.Fprint(&buf, "NoZero")
	fmt.Fprint(&buf, "() error ")
	fmt.Fprint(&buf, "{ ")
	for _, field := range fields {
		name := field.Names[0].Name

		fmt.Fprint(&buf, "if reflect.ValueOf(")
		fmt.Fprint(&buf, rcv)
		fmt.Fprint(&buf, ".")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, ").IsZero() {")
		fmt.Fprintln(&buf)
		fmt.Fprint(&buf, `	return fmt.Errorf("`)
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, ` not set")`)
		fmt.Fprint(&buf, "}")
		fmt.Fprintln(&buf)

	}
	fmt.Fprint(&buf, "	return nil")
	fmt.Fprint(&buf, "}")
	fmt.Fprintln(&buf)
	_, err := io.Copy(m.output, &buf)
	return err
}
