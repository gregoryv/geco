/*
Command gentut generates should- and must-prefixed methods handling
errors during testing.

Primary usecase is to enhance struct types that return an error with
wrapper methods that use the *testing.T Error and Fatal handlers.


*/
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/geco"
	"github.com/gregoryv/nexus"
	"github.com/gregoryv/wolf"
)

func main() {
	wolf.NewOSCmd()
	log.SetFlags(0)

	var (
		cli  = cmdline.NewParser(os.Args...)
		help = cli.Flag("-h, --help")
		pkg  = cli.Option("-p, --package").String("")
		typ  = cli.Option("-t, --type").String("")
		in   = cli.Option("-in, --input-file").String("")
		rec  = cli.Option("-rec, --receiver").String(typ + "UnderTest")
		w    = cli.Flag("-w, -write-to-out")
		out  = cli.Option("-out, --output-file").String(outFilename(in))
	)
	if rec == typ {
		// don't overwrite existing file
		out = "tut" + out
	}
	seeHelp := fmt.Sprintf(", %s --help", os.Args[0])
	switch {
	case help:
		cli.WriteUsageTo(os.Stdout)
		p, _ := nexus.NewPrinter(os.Stdout)
		p.Println("Example")
		p.Println()
		p.Println(`
//go:generate gentut --package mypkg --type Car --input-file car.go -w

same as short version

//go:generate gentut -p mypkg -t Car -in car.go -w`)
		os.Exit(0)

	case !cli.Ok():
		log.Fatal(cli.Error(), seeHelp)

	case pkg == "":
		log.Fatal("missing package", seeHelp)

	case typ == "":
		log.Fatal("missing type", seeHelp)

	case in == "":
		log.Fatal("missing input file", seeHelp)

	}

	var fh io.WriteCloser = os.Stdout
	if w {
		var err error
		fh, err = os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
		defer fh.Close()
	}
	gc := geco.NewTypeUnderTest(pkg, rec, typ, in, nil)
	if err := gc.Generate(fh); err != nil {
		log.Fatal(err)
	}
}

// outFilename returns from eg. name.go nameut_test.go
func outFilename(in string) string {
	dir := filepath.Dir(in)
	file := filepath.Base(in)
	file = strings.Replace(file, ".go", "ut_test.go", 1)
	return path.Join(dir, file)
}
