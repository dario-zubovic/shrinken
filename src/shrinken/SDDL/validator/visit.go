package validator

import "shrinken/SDDL/ast"

// this file implements Visitor pattern for AST validator

func (v *Validator) VisitPackageDef(pkg *ast.PackageDef) {

	valid := v.ValidateAttributes(pkg, pkg.AttributesList)
	if !valid {
		return
	}

	for _, attb := range pkg.AttributesList {
		attb.Accept(v)
	}
	pkg.Body.Accept(v)

}

func (v *Validator) VisitPackageBody(body *ast.PackageBody) {

	for _, importDef := range body.Imports {
		importDef.Accept(v)
	}

	for _, elem := range body.Elements {
		elem.Accept(v)
	}

}

func (v *Validator) VisitImportDef(i *ast.ImportDef) {

	valid := v.ValidateAttributes(i, i.AttributesList)
	if !valid {
		return
	}

	for _, attb := range i.AttributesList {
		attb.Accept(v)
	}

}

func (v *Validator) VisitStructDef(s *ast.StructDef) {

	valid := v.ValidateAttributes(s, s.AttributesList)
	if !valid {
		return
	}

	for _, attb := range s.AttributesList {
		attb.Accept(v)
	}
	s.Body.Accept(v)

}

func (v *Validator) VisitEnumDef(enum *ast.EnumDef) {

	valid := v.ValidateAttributes(enum, enum.AttributesList)
	if !valid {
		return
	}

	for _, attb := range enum.AttributesList {
		attb.Accept(v)
	}
	enum.Body.Accept(v)

}

func (v *Validator) VisitStructBody(structBody *ast.StructBody) {

	for _, variable := range structBody.Variables {
		variable.Accept(v)
	}

}

func (v *Validator) VisitEnumBody(enumBody *ast.EnumBody) {

	for _, e := range enumBody.Enumerals {
		e.Accept(v)
	}

}

func (v *Validator) VisitVariable(variable *ast.Variable) {
	valid := v.ValidateAttributes(variable, variable.AttributesList)
	if !valid {
		return
	}

	variable.Type.Accept(v)
	for _, attb := range variable.AttributesList {
		attb.Accept(v)
	}

}

func (v *Validator) VisitEnumeral(e *ast.Enumeral) {

}

func (v *Validator) VisitType(t *ast.Type) {
	if t.IsGeneric {

	} else if t.IsArray {

		t.ArrayChildType.Accept(v)

	} else {

	}
}

func (v *Validator) VisitAttribute(attb ast.Attribute) {

}
