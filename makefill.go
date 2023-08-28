package geco

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
)

type MakeFill struct {
	fileSet   *token.FileSet
	file      *ast.File
	typeNames []string
	output    io.Writer
}

func (m *MakeFill) Run() error {
	// find fields
	var fields []*ast.Field
	var name string

	// which structs to include
	include := make(map[string]bool)
	for _, name := range m.typeNames {
		include[name] = true
	}

	exists := func(n *ast.Field) bool {
		for _, field := range fields {
			namesEqual := field.Names[0].Name == n.Names[0].Name
			var bf bytes.Buffer
			format.Node(&bf, m.fileSet, field.Type)
			var nf bytes.Buffer
			format.Node(&nf, m.fileSet, n.Type)
			typesEqual := bf.String() == nf.String()

			if namesEqual && typesEqual {
				return true
			}
		}
		return false
	}
	ast.Inspect(m.file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			name = n.Name // save name

		case *ast.StructType:
			if include[name] {
				for _, field := range n.Fields.List {
					if !isPrivate(field.Names[0].Name) {
						continue
					}
					if exists(field) {
						continue
					}
					fields = append(fields, field)
				}
			}
		}
		return true
	})

	if len(fields) == 0 {
		return fmt.Errorf("types %v: not found", m.typeNames)
	}

	var buf bytes.Buffer
	fmt.Fprint(&buf, "func ")
	fmt.Fprint(&buf, "Fill")
	fmt.Fprint(&buf, "(")
	fmt.Fprint(&buf, "dst, src any")
	fmt.Fprint(&buf, ") ")
	fmt.Fprint(&buf, "{ ")

	for _, field := range fields {

		name = makePublic(field.Names[0].Name)

		// begin block
		fmt.Fprint(&buf, "{ ")
		fmt.Fprintln(&buf)

		// dst, dstOk := dst.(interface{SetNAME(TYPE)})
		fmt.Fprint(&buf, "dst, dstOk := dst.(interface{Set")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, "(")
		format.Node(&buf, m.fileSet, field.Type) // typeName
		fmt.Fprint(&buf, ")})")
		fmt.Fprintln(&buf)

		// src, srcOk := src.(interface{NAME() TYPE})
		fmt.Fprint(&buf, "src, srcOk := src.(interface{")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, "() ")
		format.Node(&buf, m.fileSet, field.Type) // typeName
		fmt.Fprint(&buf, "})")
		fmt.Fprintln(&buf)

		// if dstOk && srcOk {
		//     v.SetNAME(src.NAME())
		//}
		fmt.Fprint(&buf, "if dstOk && srcOk {")
		fmt.Fprintln(&buf)
		fmt.Fprint(&buf, "	dst.Set")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, "(src.")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, "()")
		fmt.Fprint(&buf, ")")
		fmt.Fprintln(&buf)
		fmt.Fprintln(&buf, "}")

		// end block
		fmt.Fprint(&buf, "}")
		fmt.Fprintln(&buf)
	}

	fmt.Fprint(&buf, "}")
	fmt.Fprintln(&buf)

	_, err := io.Copy(m.output, &buf)
	return err
}
