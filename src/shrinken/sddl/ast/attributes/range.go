package attributes

import (
	"fmt"
	"reflect"
	"shrinken/sddl/ast"
)

// RangeAttribute defines possible range of values for numeric types
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

func (attb *RangeAttribute) IsApplicable(t reflect.Type, node ast.ASTNode) (bool, error) {
	if t == reflect.TypeOf(&ast.Variable{}) {
		if node.(*ast.Variable).Type.IsGeneric {
			genericType := node.(*ast.Variable).Type.GenericType
			if genericType == ast.Float {
				return true, nil
			} else if genericType == ast.Integer32 || genericType == ast.Integer64 {
				if float64(int64(attb.Range.LowerBound)) != attb.Range.LowerBound ||
					float64(int64(attb.Range.UpperBound)) != attb.Range.UpperBound {

					return false, fmt.Errorf("Range attribute applied to integer types must be limited by integers")
				}

				return true, nil
			}
		}
	}

	return false, fmt.Errorf("Range attribute cannot be applied to non-numeric types")
}
