package ast

import (
	"shrinken/SDDL/token"
	"strconv"
)

func toStr(str interface{}) string {
	return string(str.(*token.Token).Lit)
}

func toStrUnquote(str interface{}) string {
	s := toStr(str)
	return s[1 : len(s)-2]
}

func toInt64(str interface{}) int64 {
	s := toStr(str)
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func ToFloat64(str interface{}) float64 {
	s := toStr(str)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
