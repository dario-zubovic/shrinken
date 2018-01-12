package ast

import (
	"shrinken/SDDL/token"
)

func toStr(str interface{}) string {
	return str.(*token.Token).String()
}
