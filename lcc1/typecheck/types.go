package typecheck

import (
	"lcc1/neoparser"
	"lcc1/shared"
)

type TypeCheckReturn struct {
	Expression neoparser.Expression
	Type neoparser.CompositeType
	Token shared.Token
	TokenSet *[]shared.Token
}
