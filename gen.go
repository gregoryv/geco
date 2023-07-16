// GENERATED, DO NOT EDIT!

package goref

import (
	"go/ast"
	"go/token"
	"io"
)

func (m *MakeGet) SetFileSet(v *token.FileSet) { m.fileSet = v }
func (m *MakeGet) SetFile(v *ast.File)         { m.file = v }
func (m *MakeGet) SetTypeName(v string)        { m.typeName = v }
func (m *MakeGet) SetOutput(v io.Writer)       { m.output = v }
