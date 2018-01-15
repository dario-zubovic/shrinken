package SDDL

import (
	"io/ioutil"
	"shrinken/SDDL/lexer"
	"shrinken/SDDL/parser"
	"shrinken/SDDL/validator"
	"testing"
)

func testForValidatorErrors(t *testing.T, SDDL string, expectedToBeValid bool) {
	lex := lexer.NewLexer([]byte(SDDL))
	p := parser.NewParser()

	rootNode, err := p.Parse(lex)
	if err != nil {
		t.Fatal("SDDL couldn't be parsed!", err)
		return
	}

	v := &validator.Validator{}
	valid, err := v.ValidateAST(rootNode)

	if valid != expectedToBeValid {
		if expectedToBeValid {
			t.Fatal("AST is not valid!", err)
		} else {
			t.Fatal("AST is valid, but expected it to be invalid!")
		}
	}
}

func testFileForValidatorErrors(t *testing.T, filename string, expectedToBeValid bool) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Couldn't load testing SDDL file!", err)
		return
	}

	testForValidatorErrors(t, string(b), expectedToBeValid)
}

func TestBasic(t *testing.T) {
	testFileForValidatorErrors(t, "examples/single_file/basic.sddl", true)
}

func TestUnknownType(t *testing.T) {
	testFileForValidatorErrors(t, "examples/single_file/unknown_type.sddl", false)
}
