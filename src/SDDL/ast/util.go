package ast

import "SDDL/token"

func toStr(str interface{}) string {
	return str.(*token.Token).String()
}
