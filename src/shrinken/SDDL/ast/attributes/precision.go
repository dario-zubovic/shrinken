package attributes

import (
	"fmt"
	"reflect"
	"shrinken/SDDL/ast"
)

type PrecisionAttribute struct {
	ast.Attribute
	Precision float64
}

func NewPrecisionAttribute(p interface{}) *PrecisionAttribute {
	return &PrecisionAttribute{
		Precision: p.(float64),
	}
}

func (attb *PrecisionAttribute) IsApplicable(t reflect.Type, node interface{}) (bool, error) {
	if t == reflect.TypeOf(&ast.Variable{}) {
		if node.(*ast.Variable).Type.IsGeneric &&
			node.(*ast.Variable).Type.GenericType == ast.Float {

			return true, nil
		}
	}

	if t == reflect.TypeOf(&ast.MultiVariable{}) {
		if node.(*ast.MultiVariable).Type.IsGeneric &&
			node.(*ast.MultiVariable).Type.GenericType == ast.Float {

			return true, nil
		}
	}

	return false, fmt.Errorf("Precision attribute can only be applied to float variables")
}
