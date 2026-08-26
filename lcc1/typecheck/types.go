package typecheck

import (
	"lcc1/neoparser"
	"lcc1/shared"
)

type TypeCheckReturn struct {
	Type neoparser.CompositeType
	OriginalType neoparser.CompositeType
	Token shared.Token
	TokenSet *[]shared.Token
}
