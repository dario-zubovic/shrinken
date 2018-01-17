package validator

import (
	"fmt"
	"shrinken/sddl/ast"
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

func Validate(packages []*ast.PackageDef) ([]string, error) {
	validator := &Validator{}

	for _, pkg := range packages {
		validator.traversePackage(pkg)
		if !validator.astValid {
			return nil, validator.astInvalidError
		}
	}

	err := validator.validateTypes()
	if err != nil {
		return nil, err
	}

	return validator.warnings, nil
}

// ValidateSinglePackage does single-pass AST validation
func (v *Validator) ValidateSinglePackage(pkg *ast.PackageDef) ([]string, error) {
	v.traversePackage(pkg)
	if !v.astValid {
		return nil, v.astInvalidError
	}

	err := v.validateTypes()
	if err != nil {
		return nil, err
	}

	return v.warnings, nil
}

func (v *Validator) traversePackage(pkg *ast.PackageDef) {
	v.visitor = &validatorVisitor{}
	v.visitor.Validator = v

	v.packageName = pkg.Name
	v.astValid = true

	// visit all nodes in AST, valided what is possible,
	// and collect information needed for any further validation
	pkg.Accept(v.visitor)
}

func (v *Validator) validateTypes() error {
	if len(v.declaredTypes) == 0 {
		v.addWarning("Nothing declared in package " + v.packageName + ".")
	} else {
		// check for duplicate definitions
		for i := 0; i < len(v.declaredTypes); i++ {
			t := v.declaredTypes[i]
			for n := i + 1; n < len(v.declaredTypes); n++ {
				if v.declaredTypes[n] == t {
					return fmt.Errorf("Package contains duplicate definition of type %v", t)
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
			return fmt.Errorf("Unknown type %v", t)
		}
	}

	return nil
}

func (v *Validator) addWarning(warning string) {
	v.warnings = append(v.warnings, "Warning: "+warning)
}
