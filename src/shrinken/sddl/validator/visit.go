package validator

import (
	"reflect"
	"shrinken/sddl/ast"
	"strings"
)

// this file implements Visitor pattern for AST validator

type validatorVisitor struct {
	ast.Visitor
	Validator *Validator
}

func (v *validatorVisitor) VisitPackageDef(pkg *ast.PackageDef) {

	valid := v.validateAttributes(pkg, pkg.AttributesList)
	if !valid {
		return
	}

	for _, attb := range pkg.AttributesList {
		attb.Accept(v)
	}
	pkg.Body.Accept(v)

}

func (v *validatorVisitor) VisitPackageBody(body *ast.PackageBody) {

	for _, importDef := range body.Imports {
		importDef.Accept(v)
	}

	for _, elem := range body.Elements {
		elem.Accept(v)
	}

}

func (v *validatorVisitor) VisitImportDef(i *ast.ImportDef) {

	valid := v.validateAttributes(i, i.AttributesList)
	if !valid {
		return
	}

	for _, attb := range i.AttributesList {
		attb.Accept(v)
	}

}

func (v *validatorVisitor) VisitStructDef(s *ast.StructDef) {

	valid := v.validateAttributes(s, s.AttributesList)
	if !valid {
		return
	}

	v.addDeclaredType(s.Name)

	if s.Overrides != "" {
		v.addUsedType(s.Overrides)
	}

	for _, attb := range s.AttributesList {
		attb.Accept(v)
	}
	s.Body.Accept(v)

}

func (v *validatorVisitor) VisitEnumDef(enum *ast.EnumDef) {

	valid := v.validateAttributes(enum, enum.AttributesList)
	if !valid {
		return
	}

	v.addDeclaredType(enum.Name)

	for _, attb := range enum.AttributesList {
		attb.Accept(v)
	}
	enum.Body.Accept(v)

}

func (v *validatorVisitor) VisitStructBody(structBody *ast.StructBody) {

	for _, variable := range structBody.Variables {
		variable.Accept(v)
	}

}

func (v *validatorVisitor) VisitEnumBody(enumBody *ast.EnumBody) {

	for _, e := range enumBody.Enumerals {
		e.Accept(v)
	}

}

func (v *validatorVisitor) VisitVariable(variable *ast.Variable) {
	valid := v.validateAttributes(variable, variable.AttributesList)
	if !valid {
		return
	}

	variable.Type.Accept(v)
	for _, attb := range variable.AttributesList {
		attb.Accept(v)
	}

}

func (v *validatorVisitor) VisitEnumeral(e *ast.Enumeral) {

}

func (v *validatorVisitor) VisitVariableType(t *ast.VariableType) {
	if t.IsGeneric {
		return
	}

	if t.IsArray {
		t.ArrayChildType.Accept(v)
		return
	}

	v.addUsedType(t.Name)
}

func (v *validatorVisitor) VisitAttribute(attb ast.Attribute) {

}

func (v *validatorVisitor) validateAttributes(node ast.ASTNode, attributes []ast.Attribute) bool {
	for _, attb := range attributes {
		valid, err := attb.IsApplicable(reflect.TypeOf(node), node)
		if !valid {
			v.Validator.astValid = false
			v.Validator.astInvalidError = err
			return false
		}
	}

	return true
}

func (v *validatorVisitor) addDeclaredType(name string) {
	v.Validator.declaredTypes = append(v.Validator.declaredTypes, v.Validator.packageName+"."+name)
}

func (v *validatorVisitor) addUsedType(name string) {
	if strings.Contains(name, ".") {
		v.Validator.usedTypes = append(v.Validator.usedTypes, name)
	} else {
		v.Validator.usedTypes = append(v.Validator.usedTypes, v.Validator.packageName+"."+name)
	}
}
