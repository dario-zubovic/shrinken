package ast

type PackageDef struct {
	Name string
	Body *PackageBody
}

type PackageBody struct {
	Imports  []ImportDef
	Elements []*PackageElement
}

type ImportDef string

type PackageElement struct {
}

type StructDef struct {
	PackageElement
	IsClass   bool
	Overrides string
	Name      string
	Body      *StructBody
}

type StructBody struct {
	Variables []*Variable
}

type EnumDef struct {
	PackageElement
	Name string
	Body *EnumBody
}

type EnumBody struct {
	Enumerals []*Enumeral
}

type Type struct {
	IsGeneric   bool
	GenericType GenericType

	Name string
}

type GenericType int

const (
	Integer32 GenericType = iota
	Integer64
	Short
	UnsignedInteger32
	UnsignedInteger64
	UnsignedShort
	Byte
	Bool
	String
	Char
	Float
	Double
)

type ArrayType struct {
	Type
	ChildType *Type
	Size      int // -1 to indicate that no size was specified
}

type Variable struct {
	Type *Type
	Name string
}

type Enumeral struct {
	Name string
}

func NewPackageDef(packageName string, packageBody *PackageBody) *PackageDef {
	return &PackageDef{
		Name: packageName,
		Body: packageBody,
	}
}

func NewPackageBody() *PackageBody {
	return &PackageBody{
		Imports:  make([]ImportDef, 0),
		Elements: make([]*PackageElement, 0),
	}
}

func ImportToPackageBody(body *PackageBody, importDef ImportDef) *PackageBody {
	body.Imports = append(body.Imports, importDef)
	return body
}

func AddToPackageBody(body *PackageBody, element *PackageElement) *PackageBody {
	body.Elements = append(body.Elements, element)
	return body
}

func NewImport(importName string) ImportDef {
	return ImportDef(importName)
}

func NewClassDef(name string, overrides string, body *StructBody) *StructDef {
	return &StructDef{
		IsClass:   true,
		Overrides: overrides,
		Name:      name,
		Body:      body,
	}
}

func NewStructDef(name string, overrides string, body *StructBody) *StructDef {
	return &StructDef{
		IsClass:   false,
		Overrides: overrides,
		Name:      name,
		Body:      body,
	}
}

func NewEnumDef(name string, body *EnumBody) *EnumDef {
	return &EnumDef{
		Name: name,
		Body: body,
	}
}

func NewGenericType(generic GenericType) *Type {
	return &Type{
		IsGeneric:   true,
		GenericType: generic,
	}
}

func NewType(typeName string) *Type {
	return &Type{
		IsGeneric: false,
		Name:      typeName,
	}
}

func NewArrayOfType(typeDef *Type) *ArrayType {
	return &ArrayType{
		ChildType: typeDef,
		Size:      -1,
	}
}

func NewArrayOfTypeWithSize(typeDef *Type, size int) *ArrayType {
	return &ArrayType{
		ChildType: typeDef,
		Size:      size,
	}
}

func NewVariable(typeDef *Type, name string) *Variable {
	return &Variable{
		Type: typeDef,
		Name: name,
	}
}

func NewStructBody() *StructBody {
	return &StructBody{
		Variables: make([]*Variable, 0),
	}
}

func AddToStructBody(body *StructBody, variable *Variable) *StructBody {
	body.Variables = append(body.Variables, variable)
	return body
}

func NewEnumBody() *EnumBody {
	return &EnumBody{
		Enumerals: make([]*Enumeral, 0),
	}
}

func AddToEnumBody(body *EnumBody, enumeralName string) *EnumBody {
	body.Enumerals = append(body.Enumerals, &Enumeral{
		Name: enumeralName,
	})
	return body
}
