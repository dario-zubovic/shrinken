package attributes

import (
	"fmt"
	"reflect"
	"shrinken/sddl/ast"
)

// ExportAsAttribute changes name of defined type or package in target language
type ExportAsAttribute struct {
	ast.Attribute
	ExportedName string
}

func NewExportAsAttribute(name interface{}) *ExportAsAttribute {
	return &ExportAsAttribute{
		ExportedName: ast.ToStrUnquote(name),
	}
}

func (attb *ExportAsAttribute) Accept(visitor ast.Visitor) {
	visitor.VisitAttribute(attb)
}

func (attb *ExportAsAttribute) String() string {
	return fmt.Sprint("ExportAs ", attb.ExportedName)
}

func (attb *ExportAsAttribute) IsApplicable(t reflect.Type, node ast.ASTNode) (bool, error) {
	if t == reflect.TypeOf(&ast.Variable{}) ||
		t == reflect.TypeOf(&ast.StructDef{}) ||
		t == reflect.TypeOf(&ast.EnumDef{}) ||
		t == reflect.TypeOf(&ast.PackageDef{}) {

		return true, nil
	}

	if t == reflect.TypeOf(&ast.MultiVariable{}) {
		return false, fmt.Errorf("ExportAs attribute is ambiguous between multiple variable")
	}

	return false, fmt.Errorf("ExportAs attribute can only be applied to package, classes, structs, enums or variables")
}
