package sddl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"shrinken/sddl/ast"
	"shrinken/sddl/lexer"
	"shrinken/sddl/parser"
	"shrinken/sddl/validator"
)

//go:generate ../../../bin/gocc -a SDDL.bnf

type SDDLParsed struct {
	packages []*ast.PackageDef
}

func ParseMergeAndValidate(filename string) (*SDDLParsed, error) {
	parsed, err := ParseFileOrDirectory(filename)
	if err != nil {
		return nil, err
	}

	Merge(parsed)

	warnings, err := validator.Validate(parsed.packages)
	if err != nil {
		return nil, err
	}

	if warnings != nil {
		for _, warning := range warnings {
			fmt.Println(warning)
		}
	}

	return parsed, nil
}

func ParseFileOrDirectory(filename string) (*SDDLParsed, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	parsed := &SDDLParsed{
		packages: make([]*ast.PackageDef, 0),
	}

	if stat.IsDir() {
		return parsed, parseDirectory(file, parsed)
	}

	return parsed, parseFile(filename, parsed)
}

func ParseDirectory(filename string) (*SDDLParsed, error) {
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

	parsed := &SDDLParsed{
		packages: make([]*ast.PackageDef, 0),
	}
	return parsed, parseDirectory(file, parsed)
}

func ParseFile(filename string) (*SDDLParsed, error) {
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

	parsed := &SDDLParsed{
		packages: make([]*ast.PackageDef, 0),
	}
	return parsed, parseFile(filename, parsed)
}

func parseDirectory(file *os.File, parsedList *SDDLParsed) error {
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

func parseFile(filename string, parsedList *SDDLParsed) error {
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

	parsedList.packages = append(parsedList.packages, pkg.(*ast.PackageDef))

	return nil
}
