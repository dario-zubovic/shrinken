package SDDL

import "shrinken/SDDL/ast"

func Merge(parsed *SDDLParsed) {
	if len(parsed.packages) < 2 {
		return
	}

	newPkgs := make([]*ast.PackageDef, 0)

	for _, pkg := range parsed.packages {
		merged := false
		for _, newPkg := range newPkgs {
			if newPkg.Name == pkg.Name {
				mergePackageDefs(newPkg, pkg)
				merged = true
				break
			}
		}
		if !merged {
			newPkgs = append(newPkgs, pkg)
		}
	}

	parsed.packages = newPkgs
}

func mergePackageDefs(original, other *ast.PackageDef) {
	original.AttributesList = append(original.AttributesList, other.AttributesList...)
	original.Body.Imports = append(original.Body.Imports, other.Body.Imports...)
	original.Body.Elements = append(original.Body.Elements, other.Body.Elements...)
}
