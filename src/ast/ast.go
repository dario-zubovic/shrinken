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

func NewPackageDef(packageName interface{}, packageBody interface{}) *PackageDef {
	return &PackageDef{
		Name: packageName.(string),
		Body: packageBody.(*PackageBody),
	}
}

func NewPackageBody() *PackageBody {
	return &PackageBody{
		Imports:  make([]ImportDef, 0),
		Elements: make([]*PackageElement, 0),
	}
}

func ImportToPackageBody(body interface{}, importDef interface{}) *PackageBody {
	b := body.(*PackageBody)
	b.Imports = append(b.Imports, importDef.(ImportDef))
	return b
}

func AddToPackageBody(body interface{}, element interface{}) *PackageBody {
	b := body.(*PackageBody)
	b.Elements = append(b.Elements, element.(*PackageElement))
	return b
}

func NewImport(importName interface{}) ImportDef {
	return ImportDef(importName.(string))
}

func NewClassDef(name interface{}, overrides interface{}, body interface{}) *StructDef {
	return &StructDef{
		IsClass:   true,
		Overrides: overrides.(string),
		Name:      name.(string),
		Body:      body.(*StructBody),
	}
}

func NewStructDef(name interface{}, overrides interface{}, body interface{}) *StructDef {
	return &StructDef{
		IsClass:   false,
		Overrides: overrides.(string),
		Name:      name.(string),
		Body:      body.(*StructBody),
	}
}

func NewEnumDef(name interface{}, body interface{}) *EnumDef {
	return &EnumDef{
		Name: name.(string),
		Body: body.(*EnumBody),
	}
}

func NewGenericType(generic interface{}) *Type {
	return &Type{
		IsGeneric:   true,
		GenericType: generic.(GenericType),
	}
}

func NewType(typeName interface{}) *Type {
	return &Type{
		IsGeneric: false,
		Name:      typeName.(string),
	}
}

func NewArrayOfType(typeDef interface{}) *ArrayType {
	return &ArrayType{
		ChildType: typeDef.(*Type),
		Size:      -1,
	}
}

func NewArrayOfTypeWithSize(typeDef interface{}, size interface{}) *ArrayType {
	return &ArrayType{
		ChildType: typeDef.(*Type),
		Size:      size.(int),
	}
}

func NewVariable(typeDef interface{}, name interface{}) *Variable {
	return &Variable{
		Type: typeDef.(*Type),
		Name: name.(string),
	}
}

func NewStructBody() *StructBody {
	return &StructBody{
		Variables: make([]*Variable, 0),
	}
}

func AddToStructBody(body interface{}, variable interface{}) *StructBody {
	b := body.(*StructBody)
	b.Variables = append(b.Variables, variable.(*Variable))
	return b
}

func NewEnumBody() *EnumBody {
	return &EnumBody{
		Enumerals: make([]*Enumeral, 0),
	}
}

func AddToEnumBody(body interface{}, enumeralName interface{}) *EnumBody {
	b := body.(*EnumBody)
	b.Enumerals = append(b.Enumerals, &Enumeral{
		Name: enumeralName.(string),
	})
	return b
}
