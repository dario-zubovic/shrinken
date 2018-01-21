package ast

import (
	"fmt"
	"math"
	"shrinken/sddl/token"
)

type ASTNode interface {
	Accept(visitor Visitor)
}

type PackageDef struct {
	ASTNode
	Name           string
	Body           *PackageBody
	AttributesList []Attribute
	Position       token.Pos
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
	Position       token.Pos
}

type PackageElement interface {
	ASTNode
}

type TypeDefinition interface {
	PackageElement
}

type StructDef struct {
	TypeDefinition
	IsClass          bool
	Overrides        string
	OverridesTypeDef TypeDefinition
	Name             string
	Body             *StructBody
	AttributesList   []Attribute
	Position         token.Pos
}

type StructBody struct {
	ASTNode
	Variables []*Variable
}

type EnumDef struct {
	TypeDefinition
	Name           string
	Body           *EnumBody
	AttributesList []Attribute
	Position       token.Pos
}

type EnumBody struct {
	ASTNode
	Enumerals []*Enumeral
}

type VariableType struct {
	ASTNode

	IsGeneric   bool
	GenericType GenericType

	IsArray        bool
	ArrayChildType *VariableType
	ArraySize      int // -1 to indicate that no size was specified

	Name           string
	TypeDefinition TypeDefinition
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

type Variable struct {
	ASTNode
	Type           *VariableType
	Name           string
	AttributesList []Attribute
	Position       token.Pos
}

type MultiVariable struct {
	Type           *VariableType
	Names          []string
	AttributesList []Attribute
	Positions      []token.Pos
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
		Name:     toStr(packageName),
		Body:     packageBody.(*PackageBody),
		Position: getTokenPos(packageName),
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
		Position:     getTokenPos(importName),
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
		Position:  getTokenPos(name),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewDerivedClassDef(name interface{}, overrides interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   true,
		Overrides: overrides.(string),
		Name:      toStr(name),
		Body:      body.(*StructBody),
		Position:  getTokenPos(name),
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
		Position:  getTokenPos(name),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewDerivedStructDef(name interface{}, overrides interface{}, body interface{}, attributesList interface{}) *StructDef {
	def := &StructDef{
		IsClass:   false,
		Overrides: overrides.(string),
		Name:      toStr(name),
		Body:      body.(*StructBody),
		Position:  getTokenPos(name),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewEnumDef(name interface{}, body interface{}, attributesList interface{}) *EnumDef {
	def := &EnumDef{
		Name:     toStr(name),
		Body:     body.(*EnumBody),
		Position: getTokenPos(name),
	}
	def.AttributesList = attributesList.([]Attribute)
	return def
}

func NewGenericType(generic interface{}) *VariableType {
	return &VariableType{
		IsGeneric:   true,
		GenericType: generic.(GenericType),
	}
}

func NewType(typeName interface{}) *VariableType {
	return &VariableType{
		Name: toStr(typeName),
	}
}

func NewArrayOfType(typeDef interface{}) *VariableType {
	return &VariableType{
		IsArray:        true,
		ArrayChildType: typeDef.(*VariableType),
		ArraySize:      -1,
	}
}

func NewArrayOfTypeWithSize(typeDef interface{}, size interface{}) *VariableType {
	return &VariableType{
		IsArray:        true,
		ArrayChildType: typeDef.(*VariableType),
		ArraySize:      size.(int),
	}
}

func NewTypeName(name interface{}) string {
	return toStr(name)
}

func NewVariable(typeDef interface{}, name interface{}, attributesList interface{}) *Variable {
	variable := &Variable{
		Type:     typeDef.(*VariableType),
		Name:     toStr(name),
		Position: getTokenPos(name),
	}
	variable.AttributesList = attributesList.([]Attribute)
	return variable
}

func NewMultiVariable(typeDef interface{}, firstName interface{}, secondName interface{}, attributesList interface{}) *MultiVariable {
	variable := &MultiVariable{
		Type:      typeDef.(*VariableType),
		Names:     make([]string, 2),
		Positions: make([]token.Pos, 2),
	}
	variable.Names[0] = toStr(firstName)
	variable.Positions[0] = getTokenPos(firstName)
	variable.Names[1] = toStr(secondName)
	variable.Positions[1] = getTokenPos(secondName)
	variable.AttributesList = attributesList.([]Attribute)
	return variable
}

func AddToMultiVariable(multiVariable interface{}, newName interface{}) *MultiVariable {
	multiVar := multiVariable.(*MultiVariable)
	multiVar.Names = append(multiVar.Names, toStr(newName))
	multiVar.Positions = append(multiVar.Positions, getTokenPos(newName))
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
	for i, varName := range multiVar.Names {
		b.Variables = append(b.Variables, &Variable{
			Type:           multiVar.Type,
			AttributesList: multiVar.AttributesList,
			Name:           varName,
			Position:       multiVar.Positions[i],
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
