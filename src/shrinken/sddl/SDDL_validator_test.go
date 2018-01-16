package sddl

import (
	"io/ioutil"
	"shrinken/sddl/ast"
	"shrinken/sddl/lexer"
	"shrinken/sddl/parser"
	"shrinken/sddl/validator"
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
	_, err = v.ValidateSinglePackage(rootNode.(*ast.PackageDef))

	if (err == nil) != expectedToBeValid {
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
	testFileForValidatorErrors(t, "test_data/single_file/basic.sddl", true)
}

func TestUnknownType(t *testing.T) {
	testFileForValidatorErrors(t, "test_data/single_file/unknown_type.sddl", false)
}

func TestUnknownType2(t *testing.T) {
	testFileForValidatorErrors(t, "test_data/single_file/unknown_extended_type.sddl", false)
}

func TestDuplicateDefinition(t *testing.T) {
	testFileForValidatorErrors(t, "test_data/single_file/double_definition.sddl", false)
}

func TestRangeAttribute(t *testing.T) {
	testForValidatorErrors(t, `package "test"

class Test {
	@ range: [0.14, 4] 
	int variable
}
`, false)

	testForValidatorErrors(t, `package "test"

class Test {
	@ range: [0.14, 4] 
	float variable
}
`, true)
}
