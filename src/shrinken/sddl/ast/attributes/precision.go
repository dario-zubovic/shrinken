package attributes

import (
	"fmt"
	"reflect"
	"shrinken/sddl/ast"
)

// PrecisionAttribute sets precision of floating variable
type PrecisionAttribute struct {
	ast.Attribute
	Precision float64
}

func NewPrecisionAttribute(p interface{}) *PrecisionAttribute {
	return &PrecisionAttribute{
		Precision: p.(float64),
	}
}

func (attb *PrecisionAttribute) Accept(visitor ast.Visitor) {
	visitor.VisitAttribute(attb)
}

func (attb *PrecisionAttribute) String() string {
	return fmt.Sprint("Precision ", attb.Precision)
}

func (attb *PrecisionAttribute) IsApplicable(t reflect.Type, node ast.ASTNode) (bool, error) {
	if t == reflect.TypeOf(&ast.Variable{}) {
		if node.(*ast.Variable).Type.IsGeneric &&
			node.(*ast.Variable).Type.GenericType == ast.Float {

			return true, nil
		}
	}

	return false, fmt.Errorf("Precision attribute can only be applied to float variables")
}
