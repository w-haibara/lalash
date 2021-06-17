package main

import (
	"flag"
	"os"

	"github.com/w-haibara/lalash"
)

func main() {
	c := flag.Bool("c", false, "")
	s := flag.Bool("s", false, "")
	flag.Parse()

	if *c && *s {
		panic("can not set both option -c and -s")
	}

	switch {
	case *c:
		os.Exit(lalash.RunCommand(flag.Arg(0)))
	case *s:
		os.Exit(lalash.RunScriptFile(flag.Arg(0)))
	default:
		os.Exit(lalash.RunREPL())
	}
}
