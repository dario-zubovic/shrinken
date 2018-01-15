package validator

import (
	"fmt"
	"shrinken/SDDL/ast"
)

type Validator struct {
	visitor *validatorVisitor

	packageName string

	astValid        bool
	astInvalidError error

	declaredTypes []string
	usedTypes     []string

	warnings []string
}

// ValidateAST does single-pass AST validation
func (v *Validator) ValidateAST(rootNode interface{}) (bool, error) {
	pkg, ok := rootNode.(*ast.PackageDef)
	if !ok {
		return false, fmt.Errorf("Specified root node is not package definition")
	}

	v.packageName = pkg.Name
	v.astValid = true
	v.visitor = &validatorVisitor{}
	v.visitor.Validator = v

	// visit all nodes in AST, valided what is possible,
	// and collect information needed for any further validation
	pkg.Accept(v.visitor)

	if !v.astValid {
		return false, v.astInvalidError
	}

	valid, err := v.validateTypes()
	if !valid {
		return false, err
	}

	return true, nil
}

func (v *Validator) validateTypes() (bool, error) {
	if len(v.declaredTypes) == 0 {
		v.addWarning("Nothing declared in package " + v.packageName + ".")
	} else {
		// check for duplicate definitions
		for i := 0; i < len(v.declaredTypes); i++ {
			t := v.declaredTypes[i]
			for n := i + 1; n < len(v.declaredTypes); n++ {
				if v.declaredTypes[n] == t {
					return false, fmt.Errorf("Package %v contains duplicate definition of type %v", v.packageName, t)
				}
			}
		}
	}

	for _, t := range v.usedTypes {
		exists := false
		for _, declared := range v.declaredTypes {
			if declared == t {
				exists = true
				break
			}
		}
		if !exists {
			return false, fmt.Errorf("Unknown type %v", t)
		}
	}

	return true, nil
}

func (v *Validator) addWarning(warning string) {
	v.warnings = append(v.warnings, "Warning: "+warning)
}
