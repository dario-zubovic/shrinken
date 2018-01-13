package validator

import (
	"fmt"
	"reflect"
	"shrinken/SDDL/ast"
)

type Validator struct {
	ast.Visitor

	astValid        bool
	astInvalidError error

	declaredTypes []string
	usedTypes     []string
}

// ValidateAST does single-pass AST validation
func (v *Validator) ValidateAST(rootNode interface{}) (bool, error) {
	pkg, ok := rootNode.(*ast.PackageDef)
	if !ok {
		return false, fmt.Errorf("Specified root node is not package definition")
	}

	v.astValid = true

	// visit all nodes in AST, valided what is possible,
	// and collect information needed for any further validation
	pkg.Accept(v)

	if !v.astValid {
		return false, v.astInvalidError
	}

	return true, nil
}

func (v *Validator) ValidateAttributes(node ast.ASTNode, attributes []ast.Attribute) bool {
	for _, attb := range attributes {
		valid, err := attb.IsApplicable(reflect.TypeOf(node), node)
		if !valid {
			v.astValid = false
			v.astInvalidError = err
			return false
		}
	}

	return true
}
