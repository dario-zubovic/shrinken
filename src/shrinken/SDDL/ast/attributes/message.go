package attributes

import (
	"fmt"
	"reflect"
	"shrinken/SDDL/ast"
)

type MessageAttribute struct {
	ast.Attribute
}

func NewMessageAttribute() *MessageAttribute {
	return &MessageAttribute{}
}

func (attb *MessageAttribute) Accept(visitor ast.Visitor) {
	visitor.VisitAttribute(attb)
}

func (attb *MessageAttribute) String() string {
	return "Message"
}

func (attb *MessageAttribute) IsApplicable(t reflect.Type, node ast.ASTNode) (bool, error) {
	if t == reflect.TypeOf(&ast.StructDef{}) ||
		t == reflect.TypeOf(&ast.EnumDef{}) {

		return true, nil
	}

	return false, fmt.Errorf("Message attribute is only applicable to classes, structs, and enums.")
}
