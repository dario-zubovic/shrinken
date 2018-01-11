package SDDL

import (
	"SDDL/lexer"
	"SDDL/parser"
	"testing"
)

func testForParserErrors(t *testing.T, SDDL string, valid bool) {
	lex := lexer.NewLexer([]byte(SDDL))
	p := parser.NewParser()

	_, err := p.Parse(lex)
	if (err == nil) != valid {
		if valid {
			t.Fatal("SDDL couldn't be parsed! Error: ", err)
		} else {
			t.Fatal("SDDL succesfully parsed, but expected it to fail!")
		}
	}
}

func TestPackageDecl1(t *testing.T) {
	testForParserErrors(t, `package "Test.Namespace"`, false)
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
