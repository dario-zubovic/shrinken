package analyzer

import (
	"fmt"
	"shrinken/sddl/ast"
	"shrinken/sddl/token"
	"strings"
)

// semantic analysis of parsed AST(s)
// we're basically just linking types here

type typeFinder struct {
	ast.Visitor

	packages     map[string]*ast.PackageDef
	definedTypes map[string]*definedType

	currentPackage *ast.PackageDef

	err error
}

type definedType struct {
	typeDef   ast.TypeDefinition
	parentPkg *ast.PackageDef
}

func (f *typeFinder) MapTypes(packages []*ast.PackageDef) error {
	f.packages = make(map[string]*ast.PackageDef)
	f.definedTypes = make(map[string]*definedType)

	for _, pkg := range packages {
		pkg.Accept(f)
		if f.err != nil {
			return f.err
		}
	}

	return nil
}

func (f *typeFinder) FindType(name string, parentPackage string, pos token.Pos) (*definedType, error) {
	fullName := name
	if !strings.Contains(name, ".") {
		fullName = parentPackage + "." + name
	}

	typeDef, exists := f.definedTypes[fullName]
	if !exists {
		return nil, fmt.Errorf("Unknown type %v on %v", fullName, pos)
	}

	return typeDef, nil
}

func (f *typeFinder) VisitPackageDef(pkg *ast.PackageDef) {
	_, exists := f.packages[pkg.Name]
	if exists {
		f.err = fmt.Errorf("Package %v redeclared on %v", pkg.Name, pkg.Position.String())
		return
	}

	f.packages[pkg.Name] = pkg
	f.currentPackage = pkg

	for _, attb := range pkg.AttributesList {
		attb.Accept(f)
	}
	pkg.Body.Accept(f)
}

func (f *typeFinder) VisitPackageBody(body *ast.PackageBody) {
	for _, importDef := range body.Imports {
		importDef.Accept(f)
		if f.err != nil {
			return
		}
	}

	for _, elem := range body.Elements {
		elem.Accept(f)
		if f.err != nil {
			return
		}
	}
}

func (f *typeFinder) VisitImportDef(i *ast.ImportDef) {
	for _, attb := range i.AttributesList {
		attb.Accept(f)
	}
}

func (f *typeFinder) VisitStructDef(s *ast.StructDef) {
	fullName := f.currentPackage.Name + "." + s.Name
	_, exists := f.definedTypes[fullName]
	if exists {
		var structClass string
		if s.IsClass {
			structClass = "Class"
		} else {
			structClass = "Struct"
		}
		f.err = fmt.Errorf("%v %v redeclered on %v", structClass, fullName, s.Position.String())
		return
	}
	f.definedTypes[fullName] = &definedType{
		parentPkg: f.currentPackage,
		typeDef:   s,
	}

	for _, attb := range s.AttributesList {
		attb.Accept(f)
	}
	s.Body.Accept(f)

}

func (f *typeFinder) VisitEnumDef(enum *ast.EnumDef) {
	fullName := f.currentPackage.Name + "." + enum.Name
	_, exists := f.definedTypes[fullName]
	if exists {
		f.err = fmt.Errorf("Enum %v redeclered on %v", fullName, enum.Position.String())
		return
	}
	f.definedTypes[fullName] = &definedType{
		parentPkg: f.currentPackage,
		typeDef:   enum,
	}

	for _, attb := range enum.AttributesList {
		attb.Accept(f)
	}
	enum.Body.Accept(f)

}

func (f *typeFinder) VisitStructBody(structBody *ast.StructBody) {
	for _, variable := range structBody.Variables {
		variable.Accept(f)
	}
}

func (f *typeFinder) VisitEnumBody(enumBody *ast.EnumBody) {
	for _, e := range enumBody.Enumerals {
		e.Accept(f)
	}
}

func (f *typeFinder) VisitVariable(variable *ast.Variable) {
	variable.Type.Accept(f)
	for _, attb := range variable.AttributesList {
		attb.Accept(f)
	}
}

func (f *typeFinder) VisitEnumeral(e *ast.Enumeral) {

}

func (f *typeFinder) VisitVariableType(t *ast.VariableType) {
	if t.IsGeneric {
		return
	}

	if t.IsArray {
		t.ArrayChildType.Accept(f)
		return
	}

}

func (f *typeFinder) VisitAttribute(attb ast.Attribute) {

}
