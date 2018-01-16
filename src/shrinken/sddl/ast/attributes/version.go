package attributes

// TODO: version attribute & it's handling

// I'm not really happy with handling version attribute on package level. Besides being
// easily missed, it produces unnecessary duplication of sddl files (one for each version,
// even though perhaps only one field changed in one class).
// Also, versioning of data is not specified in generation part, so speculating about
// possible attribute does more harm then help.
//                                                                 - DZubovic 16.1.2018

// import (
// 	"fmt"
// 	"reflect"
// 	"shrinken/sddl/ast"
// )

// type VersionAttribute struct {
// 	ast.Attribute
// 	Version int64
// }

// func NewVersionAttribute(version interface{}) *VersionAttribute {
// 	return &VersionAttribute{
// 		Version: ast.ToInt64(version),
// 	}
// }

// func (attb *VersionAttribute) Accept(visitor ast.Visitor) {
// 	visitor.VisitAttribute(attb)
// }

// func (attb *VersionAttribute) String() string {
// 	return fmt.Sprint("Version ", attb.Version)
// }

// func (attb *VersionAttribute) IsApplicable(t reflect.Type, node ast.ASTNode) (bool, error) {
// 	if t == reflect.TypeOf(&ast.PackageDef{}) {
// 		return true, nil
// 	}

// 	return false, fmt.Errorf("Version attribute can only be applied to package definition")
// }
