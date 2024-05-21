// GENERATED, DO NOT EDIT!

package geco

import (
	"go/ast"
	"go/token"
	"io"
)

func (m *MakeGet) SetFileSet(v *token.FileSet) { m.fileSet = v }
func (m *MakeGet) SetFile(v *ast.File)         { m.file = v }
func (m *MakeGet) SetTypeName(v string)        { m.typeName = v }
func (m *MakeGet) SetOutput(v io.Writer)       { m.output = v }
func (m *MakeSet) SetFileSet(v *token.FileSet) { m.fileSet = v }
func (m *MakeSet) SetFile(v *ast.File)         { m.file = v }
func (m *MakeSet) SetTypeName(v string)        { m.typeName = v }
func (m *MakeSet) SetOutput(v io.Writer)       { m.output = v }

func (m *MakeFill) SetFileSet(v *token.FileSet) { m.fileSet = v }
func (m *MakeFill) SetFile(v *ast.File)         { m.file = v }
func (m *MakeFill) SetTypeNames(v []string)     { m.typeNames = v }
func (m *MakeFill) SetOutput(v io.Writer)       { m.output = v }

func (m *MakeNoZero) SetFileSet(v *token.FileSet) { m.fileSet = v }
func (m *MakeNoZero) SetFile(v *ast.File)         { m.file = v }
func (m *MakeNoZero) SetTypeName(v string)        { m.typeName = v }
func (m *MakeNoZero) SetOutput(v io.Writer)       { m.output = v }

func (m *MakeEnum) SetFileSet(v *token.FileSet) { m.fileSet = v }
func (m *MakeEnum) SetFile(v *ast.File)         { m.file = v }
func (m *MakeEnum) SetTypeName(v string)        { m.typeName = v }
func (m *MakeEnum) SetOutput(v io.Writer)       { m.output = v }
