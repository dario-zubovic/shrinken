package dbgvisitor

import (
	"fmt"
	"io/ioutil"
	"os"
	"shrinken/SDDL/ast"
	"shrinken/SDDL/lexer"
	"shrinken/SDDL/parser"
)

type Visitor struct {
	ast.Visitor

	level int
}

func (v *Visitor) print(s ...interface{}) {
	indent := ""
	for i := 0; i < v.level; i++ {
		indent += " "
	}
	fmt.Println(indent, s)
}

func (v *Visitor) VisitPackageDef(pkg *ast.PackageDef) {
	print("Package:", pkg.Name)
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
	v.print("Extends:", s.Overrides)
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
		v.print("Type (generic):", t.GenericType)
	} else if t.IsArray {
		v.print("arr TODO")
		// TODO
	} else {
		v.print("Type:", t.Name)
	}
}

func (v *Visitor) VisitAttribute(attb ast.Attribute) {
	// TODO
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Please specify target SDDL file in argument!")
		return
	}

	file, err := ioutil.ReadFile(os.Args[1])
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
