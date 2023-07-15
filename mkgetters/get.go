package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func MakeGettersString(src string, typeName string) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}
	return MakeGetters(file, typeName)
}

// MakeGetters returns generated getters for the given type.
func MakeGetters(file *ast.File, typeName string) ([]byte, error) {

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
		fmt.Fprint(&buf, makePublic(name))
		fmt.Fprint(&buf, "() ")
		fmt.Fprint(&buf, field.Type.(*ast.Ident).Name)
		fmt.Fprint(&buf, " ")
		fmt.Fprint(&buf, "{ ")
		fmt.Fprint(&buf, "return ")
		fmt.Fprint(&buf, rcv)
		fmt.Fprint(&buf, ".")
		fmt.Fprint(&buf, name)
		fmt.Fprint(&buf, " ")
		fmt.Fprint(&buf, "}")
		fmt.Fprintln(&buf)
	}

	return buf.Bytes(), nil
}

func receiver(v string) string {
	return strings.ToLower(v[0:1])
}

func isPrivate(v string) bool {
	a := v[0:1]
	b := strings.ToLower(a)
	return a == b
}

func makePublic(v string) string {
	a := strings.ToUpper(v[0:1])
	return a + v[1:]
}
