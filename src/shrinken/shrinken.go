package main

import (
	"fmt"
	"shrinken/gen"
	"shrinken/sddl"
	"shrinken/sddl/dbgvisitor"

	"github.com/docopt/docopt-go"
)

var version = "v0.0.1"
var usage = `shrinken
Usage:
    shrinken <path> <output-path> <lang>
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

		r, err := sddl.ParseMergeAndAnalyze(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		dbgvisitor.PrintASTs(r)
	} else {
		path, _ := opts.String("<path>")
		outputPath, _ := opts.String("<output-path>")
		lang, _ := opts.String("<lang>")

		r, err := sddl.ParseMergeAndAnalyze(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		gen.Generate(r, lang, outputPath)
	}
}
