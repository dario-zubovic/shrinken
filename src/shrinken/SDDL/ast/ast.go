package ast

type PackageDef struct {
	Name string
	Body *PackageBody
}

type PackageBody struct {
	Imports  []*ImportDef
	Elements []PackageElement
}

type ImportDef struct {
	ImportedName   string
	AttributesList []Attribute
}

type PackageElement interface {
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
	Variables []*Variable
}

type EnumDef struct {
	PackageElement
	Name           string
	Body           *EnumBody
	AttributesList []Attribute
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
	Type           *Type
	Name           string
	AttributesList []Attribute
}

type MultiVariable struct {
	Type           *Type
	Names          []string
	AttributesList []Attribute
}

type Enumeral struct {
	Name string
}

type AttributeGroup struct {
	Body *AttributeGroupBody
}

type AttributeGroupBody struct {
	Attributes []Attribute
}

type Range interface {
}

type IntegerRange struct {
	Range
	LowerBound, UpperBound         int64
	LowerInclusive, UpperInclusive bool
}

type FloatRange struct {
	Range
	LowerBound, UpperBound         float64
	LowerInclusive, UpperInclusive bool
}

func NewPackageDef(packageName interface{}, packageBody interface{}) *PackageDef {
	return &PackageDef{
		Name: toStrUnquote(packageName),
		Body: packageBody.(*PackageBody),
	}
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
		ImportedName: toStrUnquote(importName),
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
	variable.AttributesList = attributesList.([]Attribute)
	return variable
}

func NewMultiVariable(typeDef interface{}, firstName interface{}, secondName interface{}, attributesList interface{}) *MultiVariable {
	variable := &MultiVariable{
		Type:  typeDef.(*Type),
		Names: make([]string, 2),
	}
	variable.Names[0] = toStr(firstName)
	variable.Names[1] = toStr(secondName)
	variable.AttributesList = attributesList.([]Attribute)
	return variable
}

func AddToMultiVariable(multiVariable interface{}, newName interface{}) *MultiVariable {
	multiVar := multiVariable.(*MultiVariable)
	multiVar.Names = append(multiVar.Names, toStr(newName))
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

func NewIntegerRange(lowerBound interface{}, lowerInclusive interface{}, upperBound interface{}, upperInclusive interface{}) *IntegerRange {
	return &IntegerRange{
		LowerBound:     toInt64(lowerBound),
		LowerInclusive: lowerInclusive.(bool),
		UpperBound:     toInt64(upperBound),
		UpperInclusive: upperInclusive.(bool),
	}
}

func NewFloatRange(lowerBound interface{}, lowerInclusive interface{}, upperBound interface{}, upperInclusive interface{}) *FloatRange {
	return &FloatRange{
		LowerBound:     toFloat64(lowerBound),
		LowerInclusive: lowerInclusive.(bool),
		UpperBound:     toFloat64(upperBound),
		UpperInclusive: upperInclusive.(bool),
	}
}
