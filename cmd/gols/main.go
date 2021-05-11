package main

import (
	"fmt"
	"go/ast"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/gregoryv/cmdline"
	"golang.org/x/tools/go/packages"
)

func main() {
	var (
		cli  = cmdline.NewParser(os.Args...)
		gols = Gols{
			Writer: os.Stdout,
			format: cli.Option("-f, --format").String("%n %k"),

			patterns: cli.Args(),
		}
	)
	log.SetFlags(0)
	err := gols.Run()
	if err != nil {
		log.Fatal(err)
	}

}

type Gols struct {
	io.Writer
	patterns []string
	format   string

	// state for visiting nodes
	pkg  string
	name string
	kind string
}

func (me *Gols) Run() error {

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, me.patterns...)
	if err != nil {
		return err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return fmt.Errorf("bad")
	}

	for _, pkg := range pkgs {
		me.pkg = pkg.ID
		for _, file := range pkg.Syntax {
			ast.Inspect(file, me.visit)
		}
	}
	return nil
}

func (me *Gols) visit(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.Ident:
		obj := n.Obj
		if obj == nil {
			return true
		}
		switch obj.Kind {
		case ast.Typ:
			me.name = fmt.Sprint(n.Name)
		}
	case *ast.GenDecl:
		me.kind = "struct"

	}
	me.WriteTo(me.Writer)
	return true
}

func (me *Gols) WriteTo(w io.Writer) (n int64, err error) {
	line := me.format
	if me.kind == "" || me.name == "" {
		return
	}

	line = strings.ReplaceAll(line, "%P", me.pkg)
	line = strings.ReplaceAll(line, "%p", path.Base(me.pkg))
	line = strings.ReplaceAll(line, "%n", me.name)
	line = strings.ReplaceAll(line, "%k", me.kind)
	v, err := w.Write([]byte(line + "\n"))
	// reset
	me.kind = ""
	me.name = ""
	return int64(v), err
}
