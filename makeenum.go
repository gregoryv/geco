package geco

import (
	"go/ast"
	"go/token"
	"io"
)

type MakeEnum struct {
	fileSet  *token.FileSet
	file     *ast.File
	typeName string
	output   io.Writer
	Values   []string
}

func (m *MakeEnum) Run() error {
	// find fields
	var typ string

	ast.Inspect(m.file, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.CONST {
			// We only care about const declarations.
			return true
		}

		for _, spec := range decl.Specs {
			vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
			if vspec.Type == nil && len(vspec.Values) > 0 {
				// "X = 1". With no type but a value. If the constant is untyped,
				// skip this vspec and reset the remembered type.
				typ = ""

				// If this is a simple type conversion, remember the type.
				// We don't mind if this is actually a call; a qualified call won't
				// be matched (that will be SelectorExpr, not Ident), and only unusual
				// situations will result in a function call that appears to be
				// a type conversion.
				ce, ok := vspec.Values[0].(*ast.CallExpr)
				if !ok {
					continue
				}
				id, ok := ce.Fun.(*ast.Ident)
				if !ok {
					continue
				}
				typ = id.Name
			}
			if vspec.Type != nil {
				// "X T". We have a type. Remember it.
				ident, ok := vspec.Type.(*ast.Ident)
				if !ok {
					continue
				}
				typ = ident.Name
			}
			if typ != m.typeName {
				// This is not the type we're looking for.
				continue
			}
			for _, name := range vspec.Names {
				if name.Name == "_" {
					continue
				}
				m.Values = append(m.Values, name.Name)
			}
		}
		return false
	})

	return nil
}
