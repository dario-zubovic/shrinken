package ast

import (
	"reflect"
)

type Attribute interface {
	ASTNode

	// check whether this attribute can be applied to specific node of type t
	IsApplicable(t reflect.Type, node ASTNode) (bool, error)

	// debug string describing this attribute node
	String() string
}
