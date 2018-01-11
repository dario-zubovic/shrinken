package ast

import (
	"fmt"
)

type PackageDef struct {
	Name string
	Body *PackageBody
}

type PackageBody struct {
	Imports  []*ImportDef
	Elements []*PackageElement
}

type ImportDef struct {
	Attributable
	ImportedName string
}

type PackageElement struct {
	Attributable
}

type StructDef struct {
	PackageElement
	IsClass   bool
	Overrides string // == "" if not overriding anything
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
	Attributable
	Type *Type
	Name string
}

type Enumeral struct {
	Name string
}

type Attribute struct {
	Key       string
	Value     string
	IsGroup   bool
	GroupBody *AttributeGroupBody
}

type AttributeGroupBody struct {
	Attributes []*Attribute
}

type Attributable struct {
	AttributesList []*Attribute
}

func NewPackageDef(packageName interface{}, packageBody interface{}) *PackageDef {
	return &PackageDef{
		Name: toStr(packageName),
		Body: packageBody.(*PackageBody),
	}
}

func NewPackageBody() *PackageBody {
	return &PackageBody{
		Imports:  make([]*ImportDef, 0),
		Elements: make([]*PackageElement, 0),
	}
}

func ImportToPackageBody(body interface{}, importDef interface{}) *PackageBody {
	b := body.(*PackageBody)
	b.Imports = append(b.Imports, importDef.(*ImportDef))
	return b
}

func AddToPackageBody(body interface{}, element interface{}) *PackageBody {
	b := body.(*PackageBody)
	b.Elements = append(b.Elements, element.(*PackageElement))
	return b
}

func NewImport(importName interface{}, attributesList interface{}) *ImportDef {
	def := &ImportDef{
		ImportedName: toStr(importName),
	}
	def.AttributesList = attributesList.([]*Attribute)
	return def
}

func NewClassDef(name interface{}, overrides interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   true,
		Overrides: toStr(overrides),
		Name:      toStr(name),
		Body:      body.(*StructBody),
	}
	def.AttributesList = attributesList.([]*Attribute)
	return def
}

func NewStructDef(name interface{}, overrides interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   false,
		Overrides: toStr(overrides),
		Name:      toStr(name),
		Body:      body.(*StructBody),
	}
	def.AttributesList = attributesList.([]*Attribute)
	return def
}

func NewEnumDef(name interface{}, body interface{}, attributesList interface{}) *EnumDef {
	def := &EnumDef{
		Name: toStr(name),
		Body: body.(*EnumBody),
	}
	def.AttributesList = attributesList.([]*Attribute)
	return def
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
		Name:      toStr(typeName),
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

func NewVariable(typeDef interface{}, name interface{}, attributesList interface{}) *Variable {
	variable := &Variable{
		Type: typeDef.(*Type),
		Name: toStr(name),
	}
	variable.AttributesList = attributesList.([]*Attribute)
	return variable
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
		Name: toStr(enumeralName),
	})
	return b
}

func NewKeyOnlyAttribute(key interface{}) *Attribute {
	return &Attribute{
		Key:     toStr(key),
		Value:   "",
		IsGroup: false,
	}
}

func NewAttribute(key, value interface{}) *Attribute {
	return &Attribute{
		Key:     toStr(key),
		Value:   toStr(value),
		IsGroup: false,
	}
}

func NewAttributeGroup(body interface{}) *Attribute {
	return &Attribute{
		IsGroup:   true,
		GroupBody: body.(*AttributeGroupBody),
	}
}

func NewAttributeGroupBody() *AttributeGroupBody {
	return &AttributeGroupBody{
		Attributes: make([]*Attribute, 0),
	}
}

func AddToAttributeGroupBody(body interface{}, attribute interface{}) (*AttributeGroupBody, error) {
	b := body.(*AttributeGroupBody)
	atrb := attribute.(*Attribute)
	if atrb.IsGroup {
		return nil, fmt.Errorf("didn't expect group")
	}

	b.Attributes = append(b.Attributes, atrb)
	return b, nil
}

func NewAttributesList() []*Attribute {
	return make([]*Attribute, 0)
}

func AddToAttributesList(list interface{}, attribute interface{}) []*Attribute {
	arr := list.([]*Attribute)
	arr = append(arr, attribute.(*Attribute))
	return arr
}

func AddGroupToAttributesList(list interface{}, group interface{}) ([]*Attribute, error) {
	arr := list.([]*Attribute)
	gr := group.(*Attribute)
	if !gr.IsGroup {
		return nil, fmt.Errorf("expected group")
	}

	for _, atrb := range gr.GroupBody.Attributes {
		arr = append(arr, atrb)
	}

	return arr, nil
}
