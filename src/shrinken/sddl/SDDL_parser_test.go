package sddl

import (
	"fmt"
	"io/ioutil"
	"shrinken/sddl/ast"
	"shrinken/sddl/ast/attributes"
	"shrinken/sddl/lexer"
	"shrinken/sddl/parser"
	"testing"
)

func testForParserErrors(t *testing.T, SDDL string, valid bool) {
	lex := lexer.NewLexer([]byte(SDDL))
	p := parser.NewParser()

	_, err := p.Parse(lex)
	if (err == nil) != valid {
		if valid {
			t.Fatal("SDDL couldn't be parsed!", err)
		} else {
			t.Fatal("SDDL succesfully parsed, but expected it to fail!")
		}
	}
}

func testForParserMathEvalErrors(t *testing.T, expr string, expectedResult float64) {
	SDDL := `package "dummyPackage"
@precision: %v
class DummyClass {
}`
	lex := lexer.NewLexer([]byte(fmt.Sprintf(SDDL, expr)))
	p := parser.NewParser()

	r, err := p.Parse(lex)
	if err != nil {
		t.Fatal("SDDL couldn't be parsed!", err)
		return
	}

	pkg := r.(*ast.PackageDef)
	result := pkg.Body.Elements[0].(*ast.StructDef).AttributesList[0].(*attributes.PrecisionAttribute).Precision

	if result != expectedResult {
		t.Fatalf("Wrong result; expected %v, got %v", expectedResult, result)
	}
}

func testFileForParserErrors(t *testing.T, filename string, valid bool) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Couldn't load testing SDDL file!", err)
		return
	}

	testForParserErrors(t, string(b), valid)
}

// package decl tests only check for most basic parser functions... if those fail, something is very wrong

func TestPackageDecl1(t *testing.T) {
	testForParserErrors(t, `package "Test.Namespace"`, true)
}

func TestPackageDecl2(t *testing.T) {
	testForParserErrors(t, `
		package   "Test.Namespace"

`, true)
}

func TestPackageDecl3(t *testing.T) {
	testForParserErrors(t, `package "Test_Namespace"`, true)
}

func TestPackageDecl4(t *testing.T) {
	testForParserErrors(t, `package Test.Namespace`, false)
}

func TestPackageDecl5(t *testing.T) {
	testForParserErrors(t, `package "TestNamespace`, false)
}

// Tests parser with snippet that contains all features supported by SDDL
func TestFullSnippet(t *testing.T) {
	testFileForParserErrors(t, "test_data/single_file/basic.sddl", true)
}

func TestMath1(t *testing.T) {
	testForParserMathEvalErrors(t, "42", 42)
}

func TestMath2(t *testing.T) {
	testForParserMathEvalErrors(t, "2+2", 4)
}

func TestMath3(t *testing.T) {
	testForParserMathEvalErrors(t, "2^8", 256)
}

func TestMath4(t *testing.T) {
	testForParserMathEvalErrors(t, "2^2-sqrt(8/(1+1))+5*3^2+3", 50)
}

func TestMath5(t *testing.T) {
	testForParserMathEvalErrors(t, "2.175", 2.175)
}
