package attributes

import (
	"fmt"
	"reflect"
	"shrinken/SDDL/ast"
)

type VersionAttribute struct {
	ast.Attribute
	Version int64
}

func NewVersionAttribute(version interface{}) *VersionAttribute {
	return &VersionAttribute{
		Version: ast.ToInt64(version),
	}
}

func (attb *VersionAttribute) IsApplicable(t reflect.Type, node interface{}) (bool, error) {
	if t == reflect.TypeOf(&ast.PackageDef{}) {
		return true, nil
	}

	return false, fmt.Errorf("Version attribute can only be applied to package definition")
}
