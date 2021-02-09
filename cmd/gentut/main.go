package main

import (
	"fmt"
	"log"
	"os"
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
		cli   = cmdline.NewParser(os.Args...)
		help  = cli.Flag("-h, --help")
		pkg   = cli.Option("-p, --package").String("")
		typ   = cli.Option("-t, --type").String("")
		in    = cli.Option("-in, --input-file").String("")
		tfile = "geco" + strings.ReplaceAll(in, ".go", "_test.go")
		out   = cli.Option("-out, --output-file").String(tfile)
		rec   = cli.Option("-rec, --receiver").String("tut" + typ)
	)

	seeHelp := fmt.Sprintf(", %s --help", os.Args[0])
	switch {
	case help:
		cli.WriteUsageTo(os.Stdout)
		p, _ := nexus.NewPrinter(os.Stdout)
		p.Println("Example")
		p.Println()
		p.Println("//go:generate gentut --package mypkg --type Car --input-file cars.go")
		p.Println()
		p.Println("same as short version")
		p.Println()
		p.Println("//go:generate gentut -p mypkg -t Car -in cars.go -out tutcars_test.go")

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

	gc := geco.NewTypeUnderTest(pkg, rec, typ, in, nil)
	fh, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	if err := gc.Generate(fh); err != nil {
		log.Fatal(err)
	}
}
