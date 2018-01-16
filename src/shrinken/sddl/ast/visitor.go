package ast

// this file declares Visitor interface and defines accept methods on AST nodes for double dispatching
// attribute accept methods are in shrinken/sddl/ast/attributes package

type Visitor interface {
	VisitPackageDef(pkg *PackageDef)
	VisitPackageBody(body *PackageBody)
	VisitImportDef(i *ImportDef)
	VisitStructDef(s *StructDef)
	VisitEnumDef(enum *EnumDef)
	VisitStructBody(structBody *StructBody)
	VisitEnumBody(enumBody *EnumBody)
	VisitVariable(v *Variable)
	VisitEnumeral(e *Enumeral)
	VisitType(t *Type)

	VisitAttribute(attb Attribute)
}

func (pkg *PackageDef) Accept(visitor Visitor) {
	visitor.VisitPackageDef(pkg)
}

func (body *PackageBody) Accept(visitor Visitor) {
	visitor.VisitPackageBody(body)
}

func (i *ImportDef) Accept(visitor Visitor) {
	visitor.VisitImportDef(i)
}

func (s *StructDef) Accept(visitor Visitor) {
	visitor.VisitStructDef(s)
}

func (enum *EnumDef) Accept(visitor Visitor) {
	visitor.VisitEnumDef(enum)
}

func (structBody *StructBody) Accept(visitor Visitor) {
	visitor.VisitStructBody(structBody)
}

func (enumBody *EnumBody) Accept(visitor Visitor) {
	visitor.VisitEnumBody(enumBody)
}

func (t *Type) Accept(visitor Visitor) {
	visitor.VisitType(t)
}

func (v *Variable) Accept(visitor Visitor) {
	visitor.VisitVariable(v)
}

func (e *Enumeral) Accept(visitor Visitor) {
	visitor.VisitEnumeral(e)
}
