package attributes

import (
	"fmt"
	"reflect"
	"shrinken/SDDL/ast"
)

type RangeAttribute struct {
	ast.Attribute
	Range *ast.Range
}

func NewRangeAttribute(r interface{}) *RangeAttribute {
	return &RangeAttribute{
		Range: r.(*ast.Range),
	}
}

func (attb *RangeAttribute) Accept(visitor ast.Visitor) {
	visitor.VisitAttribute(attb)
}

func (attb *RangeAttribute) String() string {
	return fmt.Sprint("Range ", ast.RangeToString(attb.Range))
}

func (attb *RangeAttribute) IsApplicable(t reflect.Type, node interface{}) (bool, error) {
	if t == reflect.TypeOf(&ast.Variable{}) {
		if node.(*ast.Variable).Type.IsGeneric {
			genericType := node.(*ast.Variable).Type.GenericType
			if genericType == ast.Integer32 || genericType == ast.Integer64 || genericType == ast.Float {
				return true, nil
			}
		}
	}

	if t == reflect.TypeOf(&ast.MultiVariable{}) {
		if node.(*ast.MultiVariable).Type.IsGeneric {
			genericType := node.(*ast.MultiVariable).Type.GenericType
			if genericType == ast.Integer32 || genericType == ast.Integer64 || genericType == ast.Float {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("Range attribute cannot be applied to non-numeric types")
}
