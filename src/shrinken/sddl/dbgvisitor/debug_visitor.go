package dbgvisitor

import (
	"fmt"
	"io/ioutil"
	"os"
	"shrinken/sddl/ast"
	"shrinken/sddl/lexer"
	"shrinken/sddl/parser"
	"strconv"
)

type Visitor struct {
	ast.Visitor

	level int
}

func (v *Visitor) print(s ...interface{}) {
	str := ""
	for i := 0; i < v.level; i++ {
		str += "    "
	}
	for i, p := range s {
		if i == 0 {
			str += fmt.Sprint(p)
		} else {
			str += " " + fmt.Sprint(p)
		}
	}
	fmt.Println(str)
}

func (v *Visitor) VisitPackageDef(pkg *ast.PackageDef) {
	v.print("Package:", pkg.Name)
	v.level++
	for _, attb := range pkg.AttributesList {
		attb.Accept(v)
	}
	pkg.Body.Accept(v)
	v.level--
}

func (v *Visitor) VisitPackageBody(body *ast.PackageBody) {
	v.print("{")
	v.level++
	for _, importDef := range body.Imports {
		importDef.Accept(v)
	}
	for _, elem := range body.Elements {
		elem.Accept(v)
	}
	v.level--
	v.print("}")
}

func (v *Visitor) VisitImportDef(i *ast.ImportDef) {
	v.print("Import:", i.ImportedName)

	v.level++
	for _, attb := range i.AttributesList {
		attb.Accept(v)
	}
	v.level--
}

func (v *Visitor) VisitStructDef(s *ast.StructDef) {
	v.print("Struct:", s.Name)
	v.level++
	if s.Overrides != "" {
		v.print("Extends:", s.Overrides)
	}
	v.print("Class:", s.IsClass)
	for _, attb := range s.AttributesList {
		attb.Accept(v)
	}
	s.Body.Accept(v)
	v.level--
}

func (v *Visitor) VisitEnumDef(enum *ast.EnumDef) {
	v.print("Enum:", enum.Name)
	v.level++
	for _, attb := range enum.AttributesList {
		attb.Accept(v)
	}
	enum.Body.Accept(v)
	v.level--
}

func (v *Visitor) VisitStructBody(structBody *ast.StructBody) {
	v.print("{")
	v.level++
	for _, variable := range structBody.Variables {
		variable.Accept(v)
	}
	v.level--
	v.print("}")
}

func (v *Visitor) VisitEnumBody(enumBody *ast.EnumBody) {
	v.print("{")
	v.level++
	for _, e := range enumBody.Enumerals {
		e.Accept(v)
	}
	v.level--
	v.print("}")
}

func (v *Visitor) VisitVariable(variable *ast.Variable) {
	v.print("Variable:", variable.Name)
	v.level++
	variable.Type.Accept(v)
	for _, attb := range variable.AttributesList {
		attb.Accept(v)
	}
	v.level--
}

func (v *Visitor) VisitEnumeral(e *ast.Enumeral) {
	v.print("Enumeral:", e.Name)
}

func (v *Visitor) VisitType(t *ast.Type) {
	if t.IsGeneric {
		v.print("Type (generic):", t.GenericType.String())
	} else if t.IsArray {
		size := ""
		if t.ArraySize == -1 {
			size = "∞"
		} else {
			size = strconv.Itoa(t.ArraySize)
		}
		v.print("Array of size", size, ". Child type:")
		v.level++
		t.ArrayChildType.Accept(v)
		v.level--
	} else {
		v.print("Type:", t.Name)
	}
}

func (v *Visitor) VisitAttribute(attb ast.Attribute) {
	v.print("Attribute:", attb.String())
}

func PrintAST(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading SDDL file!", err)
		return
	}

	p := parser.NewParser()
	lex := lexer.NewLexer(file)

	r, err := p.Parse(lex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing SDDL file!", err)
		return
	}

	debugVisitor := &Visitor{}

	pkg := r.(*ast.PackageDef)
	pkg.Accept(debugVisitor)
}
