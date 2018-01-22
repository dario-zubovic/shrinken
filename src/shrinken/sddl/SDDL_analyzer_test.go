package sddl

import (
	"io/ioutil"
	"shrinken/sddl/analyzer"
	"shrinken/sddl/ast"
	"shrinken/sddl/lexer"
	"shrinken/sddl/parser"
	"testing"
)

func testForAnalyzerErrors(t *testing.T, SDDL string, expectedToBeValid bool) {
	lex := lexer.NewLexer([]byte(SDDL))
	p := parser.NewParser()

	rootNode, err := p.Parse(lex)
	if err != nil {
		t.Fatal("SDDL couldn't be parsed!", err)
		return
	}

	pkgs := make([]*ast.PackageDef, 1)
	pkgs[0] = rootNode.(*ast.PackageDef)
	err = analyzer.Analyze(pkgs)

	if (err == nil) != expectedToBeValid {
		if expectedToBeValid {
			t.Fatal("AST is not valid!", err)
		} else {
			t.Fatal("AST is valid, but expected it to be invalid!")
		}
	}
}

func testFileForAnalyzerErrors(t *testing.T, filename string, expectedToBeValid bool) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Couldn't load testing SDDL file!", err)
		return
	}

	testForAnalyzerErrors(t, string(b), expectedToBeValid)
}

func testFolderForAnalyzerErrors(t *testing.T, filename string, expectedToBeValid bool) {
	_, err := ParseMergeAndAnalyze(filename)
	if (err == nil) != expectedToBeValid {
		if expectedToBeValid {
			t.Fatal("AST is not valid!", err)
		} else {
			t.Fatal("AST is valid, but expected it to be invalid!")
		}
	}
}

func TestBasic(t *testing.T) {
	testFileForAnalyzerErrors(t, "test_data/single_file/basic.sddl", true)
}

func TestUnknownType(t *testing.T) {
	testFileForAnalyzerErrors(t, "test_data/single_file/unknown_type.sddl", false)
}

func TestUnknownType2(t *testing.T) {
	testFileForAnalyzerErrors(t, "test_data/single_file/unknown_extended_type.sddl", false)
}

func TestDuplicateDefinition(t *testing.T) {
	testFileForAnalyzerErrors(t, "test_data/single_file/double_definition.sddl", false)
}

func TestVariableHiding(t *testing.T) {
	testFileForAnalyzerErrors(t, "test_data/single_file/variable_hiding.sddl", false)
}

func TestRangeAttribute(t *testing.T) {
	testForAnalyzerErrors(t, `package test

class Test {
	@ range: [0.14, 4] 
	int variable
}
`, false)

	testForAnalyzerErrors(t, `package test

class Test {
	@ range: [0.14, 4] 
	float variable
}
`, true)
}

func TestSameNameClasses(t *testing.T) {
	testFolderForAnalyzerErrors(t, "test_data/multipkg/same_name/", true)
}
