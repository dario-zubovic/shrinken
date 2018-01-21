package ast

import (
	"shrinken/sddl/token"
	"strconv"
)

func toStr(str interface{}) string {
	return string(str.(*token.Token).Lit)
}

func getTokenPos(tok interface{}) token.Pos {
	return tok.(*token.Token).Pos
}

func ToStrUnquote(str interface{}) string {
	s := toStr(str)
	return s[1 : len(s)-1]
}

func ToInt64(str interface{}) int64 {
	s := toStr(str)
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func ToFloat64(str interface{}) float64 {
	s := toStr(str)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func RangeToString(r *Range) string {
	s := ""
	if r.LowerInclusive {
		s += "["
	} else {
		s += "<"
	}

	s += strconv.FormatFloat(r.LowerBound, 'g', -1, 64) + ", " + strconv.FormatFloat(r.UpperBound, 'g', -1, 64)

	if r.UpperInclusive {
		s += "]"
	} else {
		s += ">"
	}

	return s
}
