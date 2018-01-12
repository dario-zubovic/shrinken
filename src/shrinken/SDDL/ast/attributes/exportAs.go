package attributes

import (
	"fmt"
	"reflect"
	"shrinken/SDDL/ast"
)

type ExportAsAttribute struct {
	ast.Attribute
	ExportedName string
}

func NewExportAsAttribute(name interface{}) *ExportAsAttribute {
	return &ExportAsAttribute{
		ExportedName: ast.ToStrUnquote(name),
	}
}

func (attb *ExportAsAttribute) IsApplicable(t reflect.Type, node interface{}) (bool, error) {
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
