package ast

import (
	"shrinken/SDDL/token"
	"strconv"
)

func toStr(str interface{}) string {
	return str.(*token.Token).String()
}

func toStrUnquote(str interface{}) string {
	s := str.(*token.Token).String()
	return s[1 : len(s)-2]
}

func toInt64(str interface{}) int64 {
	s := toStr(str)
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func toFloat64(str interface{}) float64 {
	s := toStr(str)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
