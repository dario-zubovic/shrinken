package sddl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"shrinken/sddl/analyzer"
	"shrinken/sddl/ast"
	"shrinken/sddl/lexer"
	"shrinken/sddl/parser"
	"shrinken/sddl/validator"
)

//go:generate ../../../bin/gocc -a SDDL.bnf

type SDDLTree struct {
	Packages []*ast.PackageDef
}

func ParseMergeAndAnalyze(filename string) (*SDDLTree, error) {
	tree, err := ParseFileOrDirectory(filename)
	if err != nil {
		return nil, err
	}

	Merge(tree)

	err = analyzer.Analyze(tree.Packages)
	if err != nil {
		return nil, err
	}

	warnings, err := validator.Validate(tree.Packages)
	if err != nil {
		return nil, err
	}

	if warnings != nil {
		for _, warning := range warnings {
			fmt.Println(warning)
		}
	}

	return tree, nil
}

func ParseFileOrDirectory(filename string) (*SDDLTree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	parsed := &SDDLTree{
		Packages: make([]*ast.PackageDef, 0),
	}

	if stat.IsDir() {
		return parsed, parseDirectory(file, parsed)
	}

	return parsed, parseFile(filename, parsed)
}

func ParseDirectory(filename string) (*SDDLTree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("ParseDirectory is only accepting directories")
	}

	parsed := &SDDLTree{
		Packages: make([]*ast.PackageDef, 0),
	}
	return parsed, parseDirectory(file, parsed)
}

func ParseFile(filename string) (*SDDLTree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("ParseFile is only accepting files")
	}

	parsed := &SDDLTree{
		Packages: make([]*ast.PackageDef, 0),
	}
	return parsed, parseFile(filename, parsed)
}

func parseDirectory(file *os.File, parsedList *SDDLTree) error {
	infos, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			childFile, err := os.Open(info.Name())
			if err != nil {
				return err
			}
			defer file.Close()

			err = parseDirectory(childFile, parsedList)
			if err != nil {
				return err
			}
		} else {
			err = parseFile(filepath.Join(file.Name(), info.Name()), parsedList)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parseFile(filename string, parsedList *SDDLTree) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lex := lexer.NewLexer(file)
	p := parser.NewParser()

	pkg, err := p.Parse(lex)
	if err != nil {
		return err
	}

	parsedList.Packages = append(parsedList.Packages, pkg.(*ast.PackageDef))

	return nil
}
