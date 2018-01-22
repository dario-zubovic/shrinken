package analyzer

import (
	"fmt"
	"reflect"
	"shrinken/sddl/ast"
	"shrinken/sddl/token"
)

// semantic analysis of parsed AST(s)
// we're basically just linking types, type checking,
// inheritance checking and validating attributes

type StaticAnalyzer struct {
	ast.Visitor

	currentPkg  *ast.PackageDef
	variablePos token.Pos
	err         error

	finder *typeFinder
}

func Analyze(packages []*ast.PackageDef) error {
	analyzer := &StaticAnalyzer{}
	return analyzer.Analyze(packages)
}

func (a *StaticAnalyzer) Analyze(packages []*ast.PackageDef) error {
	// first we're letting typeFinder to go through AST of each package
	// to make a map of all type definitions

	a.finder = &typeFinder{}
	a.finder.MapTypes(packages)

	// next we're linking all variable types to their definitions
	// and doing any validations that are left

	for _, pkg := range packages {
		pkg.Accept(a)
		if a.err != nil {
			return a.err
		}
	}

	return nil
}

func (a *StaticAnalyzer) VisitPackageDef(pkg *ast.PackageDef) {
	a.currentPkg = pkg

	a.validateAttributes(pkg, pkg.AttributesList)
	if a.err != nil {
		return
	}

	for _, attb := range pkg.AttributesList {
		attb.Accept(a)
	}
	pkg.Body.Accept(a)
}

func (a *StaticAnalyzer) VisitPackageBody(body *ast.PackageBody) {
	for _, importDef := range body.Imports {
		importDef.Accept(a)
		if a.err != nil {
			return
		}
	}

	for _, elem := range body.Elements {
		elem.Accept(a)
		if a.err != nil {
			return
		}
	}
}

func (a *StaticAnalyzer) VisitImportDef(i *ast.ImportDef) {
	a.validateAttributes(i, i.AttributesList)
	if a.err != nil {
		return
	}

	for _, attb := range i.AttributesList {
		attb.Accept(a)
	}
}

func (a *StaticAnalyzer) VisitStructDef(s *ast.StructDef) {
	if s.Overrides != "" {
		def, err := a.finder.FindType(s.Overrides, a.currentPkg.Name, s.Position)
		if err != nil {
			a.err = err
			return
		}

		parent := &definedType{
			parentPkg: a.currentPkg,
			typeDef:   s,
		}
		chain := make([]*ast.StructDef, 1)
		chain[0] = s
		err = a.checkStructInheritance(parent, def, chain)
		if err != nil {
			a.err = err
			return
		}

		s.OverridesTypeDef = def.typeDef
	}

	a.validateAttributes(s, s.AttributesList)
	if a.err != nil {
		return
	}

	for _, attb := range s.AttributesList {
		attb.Accept(a)
	}

	s.Body.Accept(a)
}

func (a *StaticAnalyzer) VisitEnumDef(enum *ast.EnumDef) {
	a.validateAttributes(enum, enum.AttributesList)
	if a.err != nil {
		return
	}

	for _, attb := range enum.AttributesList {
		attb.Accept(a)
	}

	enum.Body.Accept(a)
}

func (a *StaticAnalyzer) VisitStructBody(structBody *ast.StructBody) {
	for _, variable := range structBody.Variables {
		variable.Accept(a)
		if a.err != nil {
			return
		}
	}
}

func (a *StaticAnalyzer) VisitEnumBody(enumBody *ast.EnumBody) {
	for _, e := range enumBody.Enumerals {
		e.Accept(a)
	}
}

func (a *StaticAnalyzer) VisitVariable(variable *ast.Variable) {
	a.variablePos = variable.Position
	variable.Type.Accept(a)
	if a.err != nil {
		return
	}

	a.validateAttributes(variable, variable.AttributesList)
	if a.err != nil {
		return
	}

	for _, attb := range variable.AttributesList {
		attb.Accept(a)
	}
}

func (a *StaticAnalyzer) VisitEnumeral(e *ast.Enumeral) {

}

func (a *StaticAnalyzer) VisitVariableType(t *ast.VariableType) {
	if t.IsGeneric {
		return
	}

	if t.IsArray {
		t.ArrayChildType.Accept(a)
		return
	}

	def, err := a.finder.FindType(t.Name, a.currentPkg.Name, a.variablePos)
	if err != nil {
		a.err = err
		return
	}
	t.TypeDefinition = def.typeDef
}

func (a *StaticAnalyzer) VisitAttribute(attb ast.Attribute) {

}

func (a *StaticAnalyzer) checkStructInheritance(parent, child *definedType, chain []*ast.StructDef) error {
	// in context of checkStructInheritance parent is struct that inherits from child - don't ask why :)

	parentStruct := parent.typeDef.(*ast.StructDef)

	// check if child is indeed *ast.StructDef
	if reflect.TypeOf(child.typeDef) != reflect.TypeOf(parentStruct) {
		return fmt.Errorf("Struct or class %v extends type which is not a struct or class", parentStruct.Name)
	}

	childStruct := child.typeDef.(*ast.StructDef)

	// make sure that both parent and child are either class or struct
	if parentStruct.IsClass != childStruct.IsClass {
		if parentStruct.IsClass {
			return fmt.Errorf("Class %v cannot extend struct %v on %v", parentStruct.Name, childStruct.Name, parentStruct.Position.String())
		}
		return fmt.Errorf("Struct %v cannot extend class %v on %v", parentStruct.Name, childStruct.Name, parentStruct.Position.String())
	}

	// check for circular inheritance
	chainStr := ""
	for _, c := range chain {
		chainStr += c.Name + ", "
		if c == childStruct {
			return fmt.Errorf("Circular inheritance detected (%v)", chainStr)
		}
	}

	if childStruct.Overrides == "" {
		return nil
	}

	// check for circular inheritance recursively for all child structs

	chain = append(chain, childStruct)

	childsChild, err := a.finder.FindType(childStruct.Overrides, child.parentPkg.Name, childStruct.Position)
	if err != nil {
		return err
	}

	err = a.checkStructInheritance(child, childsChild, chain)
	if err != nil {
		return err
	}

	return nil
}

func (a *StaticAnalyzer) validateAttributes(node ast.ASTNode, attributes []ast.Attribute) {
	for _, attb := range attributes {
		valid, err := attb.IsApplicable(reflect.TypeOf(node), node)
		if !valid {
			a.err = err
			return
		}
	}
}
