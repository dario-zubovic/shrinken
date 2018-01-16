package main

import (
	"fmt"
	"shrinken/sddl/dbgvisitor"

	"github.com/docopt/docopt-go"
)

var version = "v0.0.1"
var usage = `shrinken
Usage:
  shrinken <path>
  shrinken print-ast <path>
  shrinken (-h | --help)
  shrinken --version

Options:
  -h --help    Show this screen.
  --version    Show version.

Commands:
  print-ast    Print parsed AST from specified file. Only for debugging.
`

func main() {

	opts, err := docopt.ParseArgs(usage, nil, version)
	if err != nil {
		fmt.Println("Error parsing arguments:", err)
		return
	}

	if opts["print-ast"] != nil {
		path, _ := opts.String("<path>")
		dbgvisitor.PrintAST(path)
	} else {
		// TODO: main entry
	}
}
