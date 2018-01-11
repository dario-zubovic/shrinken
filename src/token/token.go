// Code generated by gocc; DO NOT EDIT.

package token

import (
	"fmt"
)

type Token struct {
	Type
	Lit []byte
	Pos
}

type Type int

const (
	INVALID Type = iota
	EOF
)

type Pos struct {
	Offset int
	Line   int
	Column int
}

func (p Pos) String() string {
	return fmt.Sprintf("Pos(offset=%d, line=%d, column=%d)", p.Offset, p.Line, p.Column)
}

type TokenMap struct {
	typeMap []string
	idMap   map[string]Type
}

func (m TokenMap) Id(tok Type) string {
	if int(tok) < len(m.typeMap) {
		return m.typeMap[tok]
	}
	return "unknown"
}

func (m TokenMap) Type(tok string) Type {
	if typ, exist := m.idMap[tok]; exist {
		return typ
	}
	return INVALID
}

func (m TokenMap) TokenString(tok *Token) string {
	//TODO: refactor to print pos & token string properly
	return fmt.Sprintf("%s(%d,%s)", m.Id(tok.Type), tok.Type, tok.Lit)
}

func (m TokenMap) StringType(typ Type) string {
	return fmt.Sprintf("%s(%d)", m.Id(typ), typ)
}

var TokMap = TokenMap{
	typeMap: []string{
		"INVALID",
		"$",
		"package",
		"str",
		"empty",
		"use",
		"class",
		"letters",
		"{",
		"}",
		":",
		"struct",
		"int",
		"int32",
		"int64",
		"long",
		"bool",
		"short",
		"uint",
		"uint32",
		"uint64",
		"ulong",
		"ubool",
		"ushort",
		"byte",
		"string",
		"char",
		"float",
		"double",
		"[]",
		"[",
		"integer",
		"]",
	},

	idMap: map[string]Type{
		"INVALID": 0,
		"$":       1,
		"package": 2,
		"str":     3,
		"empty":   4,
		"use":     5,
		"class":   6,
		"letters": 7,
		"{":       8,
		"}":       9,
		":":       10,
		"struct":  11,
		"int":     12,
		"int32":   13,
		"int64":   14,
		"long":    15,
		"bool":    16,
		"short":   17,
		"uint":    18,
		"uint32":  19,
		"uint64":  20,
		"ulong":   21,
		"ubool":   22,
		"ushort":  23,
		"byte":    24,
		"string":  25,
		"char":    26,
		"float":   27,
		"double":  28,
		"[]":      29,
		"[":       30,
		"integer": 31,
		"]":       32,
	},
}
