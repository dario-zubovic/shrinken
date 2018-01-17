package ast

import (
	"fmt"
	"math"
)

type ASTNode interface {
	Accept(visitor Visitor)
}

type PackageDef struct {
	ASTNode
	Name           string
	Body           *PackageBody
	AttributesList []Attribute
}

type PackageBody struct {
	ASTNode
	Imports  []*ImportDef
	Elements []PackageElement
}

type ImportDef struct {
	ASTNode
	ImportedName   string
	AttributesList []Attribute
}

type PackageElement interface {
	ASTNode
}

type StructDef struct {
	PackageElement
	IsClass        bool
	Overrides      string // == "" if not overriding anything
	Name           string
	Body           *StructBody
	AttributesList []Attribute
}

type StructBody struct {
	ASTNode
	Variables []*Variable
}

type EnumDef struct {
	PackageElement
	Name           string
	Body           *EnumBody
	AttributesList []Attribute
}

type EnumBody struct {
	ASTNode
	Enumerals []*Enumeral
}

type Type struct {
	ASTNode

	IsGeneric   bool
	GenericType GenericType

	IsArray        bool
	ArrayChildType *Type
	ArraySize      int // -1 to indicate that no size was specified

	Name string
}

//go:generate stringer -type=GenericType
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

type VariableName struct {
	Name    string
	Package string
}

type Variable struct {
	ASTNode
	Type           *Type
	Name           *VariableName
	AttributesList []Attribute
}

type MultiVariable struct {
	Type           *Type
	Names          []*VariableName
	AttributesList []Attribute
}

type Enumeral struct {
	ASTNode
	Name string
}

type AttributeGroup struct {
	Body *AttributeGroupBody
}

type AttributeGroupBody struct {
	Attributes []Attribute
}

type Range struct {
	LowerBound, UpperBound         float64
	LowerInclusive, UpperInclusive bool
}

func NewPackageDef(packageName interface{}, packageBody interface{}, attributesList interface{}) *PackageDef {
	def := &PackageDef{
		Name: toStr(packageName),
		Body: packageBody.(*PackageBody),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewPackageBody() *PackageBody {
	return &PackageBody{
		Imports:  make([]*ImportDef, 0),
		Elements: make([]PackageElement, 0),
	}
}

func ImportToPackageBody(body interface{}, importDef interface{}) *PackageBody {
	b := body.(*PackageBody)
	b.Imports = append(b.Imports, importDef.(*ImportDef))
	return b
}

func AddToPackageBody(body interface{}, element interface{}) *PackageBody {
	b := body.(*PackageBody)
	b.Elements = append(b.Elements, element.(PackageElement))
	return b
}

func NewImport(importName interface{}, attributesList interface{}) *ImportDef {
	def := &ImportDef{
		ImportedName: ToStrUnquote(importName),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewClassDef(name interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   true,
		Overrides: "",
		Name:      toStr(name),
		Body:      body.(*StructBody),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewDerivedClassDef(name interface{}, overrides interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   true,
		Overrides: toStr(overrides),
		Name:      toStr(name),
		Body:      body.(*StructBody),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewStructDef(name interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   false,
		Overrides: "",
		Name:      toStr(name),
		Body:      body.(*StructBody),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewDerivedStructDef(name interface{}, overrides interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   false,
		Overrides: toStr(overrides),
		Name:      toStr(name),
		Body:      body.(*StructBody),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewEnumDef(name interface{}, body interface{}, attributesList interface{}) *EnumDef {
	def := &EnumDef{
		Name: toStr(name),
		Body: body.(*EnumBody),
	}
	def.AttributesList = attributesList.([]Attribute)
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
		Name: toStr(typeName),
	}
}

func NewArrayOfType(typeDef interface{}) *Type {
	return &Type{
		IsArray:        true,
		ArrayChildType: typeDef.(*Type),
		ArraySize:      -1,
	}
}

func NewArrayOfTypeWithSize(typeDef interface{}, size interface{}) *Type {
	return &Type{
		IsArray:        true,
		ArrayChildType: typeDef.(*Type),
		ArraySize:      size.(int),
	}
}

func NewVariableName(name interface{}, packageName interface{}) *VariableName {
	varName := &VariableName{
		Name: toStr(name),
	}
	if packageName != nil {
		varName.Package = toStr(packageName)
	}

	return varName
}

func NewVariable(typeDef interface{}, name interface{}, attributesList interface{}) *Variable {
	variable := &Variable{
		Type: typeDef.(*Type),
		Name: name.(*VariableName),
	}
	variable.AttributesList = attributesList.([]Attribute)
	return variable
}

func NewMultiVariable(typeDef interface{}, firstName interface{}, secondName interface{}, attributesList interface{}) *MultiVariable {
	variable := &MultiVariable{
		Type:  typeDef.(*Type),
		Names: make([]*VariableName, 2),
	}
	variable.Names[0] = firstName.(*VariableName)
	variable.Names[1] = secondName.(*VariableName)
	variable.AttributesList = attributesList.([]Attribute)
	return variable
}

func AddToMultiVariable(multiVariable interface{}, newName interface{}) *MultiVariable {
	multiVar := multiVariable.(*MultiVariable)
	multiVar.Names = append(multiVar.Names, newName.(*VariableName))
	return multiVar
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

func AddMultiVariableToStructBody(body interface{}, multiVariable interface{}) *StructBody {
	b := body.(*StructBody)
	multiVar := multiVariable.(*MultiVariable)
	for _, varName := range multiVar.Names {
		b.Variables = append(b.Variables, &Variable{
			Type:           multiVar.Type,
			AttributesList: multiVar.AttributesList,
			Name:           varName,
		})
	}
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

func NewAttributeGroup(body interface{}) *AttributeGroup {
	return &AttributeGroup{
		Body: body.(*AttributeGroupBody),
	}
}

func NewAttributeGroupBody() *AttributeGroupBody {
	return &AttributeGroupBody{
		Attributes: make([]Attribute, 0),
	}
}

func AddToAttributeGroupBody(body interface{}, attribute interface{}) *AttributeGroupBody {
	b := body.(*AttributeGroupBody)
	atrb := attribute.(Attribute)
	b.Attributes = append(b.Attributes, atrb)
	return b
}

func NewAttributesList() []Attribute {
	return make([]Attribute, 0)
}

func AddToAttributesList(list interface{}, attribute interface{}) []Attribute {
	arr := list.([]Attribute)
	arr = append(arr, attribute.(Attribute))
	return arr
}

func AddGroupToAttributesList(list interface{}, group interface{}) []Attribute {
	arr := list.([]Attribute)
	gr := group.(*AttributeGroup)

	for _, atrb := range gr.Body.Attributes {
		arr = append(arr, atrb)
	}

	return arr
}

func NewRange(lowerBound interface{}, lowerInclusive interface{}, upperBound interface{}, upperInclusive interface{}) (*Range, error) {
	r := &Range{
		LowerBound:     lowerBound.(float64),
		LowerInclusive: lowerInclusive.(bool),
		UpperBound:     upperBound.(float64),
		UpperInclusive: upperInclusive.(bool),
	}

	if r.LowerBound > r.UpperBound {
		return nil, fmt.Errorf("Lower bound on range is higher than upper bound")
	}

	if math.IsInf(r.LowerBound, 0) && r.LowerInclusive ||
		math.IsInf(r.UpperBound, 0) && r.UpperInclusive {

		return nil, fmt.Errorf("Infinity cannot be inclusive in range")
	}

	return r, nil
}
