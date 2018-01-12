package ast

import (
	"reflect"
)

type Attribute interface {
	IsApplicable(t reflect.Type, node interface{}) (bool, error)
}
