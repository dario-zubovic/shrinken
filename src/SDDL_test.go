package SDDL

import (
	"lexer"
	"parser"
	"testing"
)

var testData = `
package "Test.Namespace"

use "Other.Namespace"

`

func testForParserErrors(t *testing.T, SDDL string, name string, valid bool) {
	lex := lexer.NewLexer([]byte(SDDL))
	p := parser.NewParser()

	result, err := p.Parse(lex)
	if (err == nil) != valid {
		if valid {
			t.Fatal("SDDL '", name, "' couldn't be parsed! Error: ", err)
		} else {
			t.Fatal("SDDL '", name, "' succesfully parsed, but expected it to fail!")
		}
	}
}

func TestPackageDecl(t *testing.T) {
	testForParserErrors(t, testData, "package decl", true)
}
