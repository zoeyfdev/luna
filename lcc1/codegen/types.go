package codegen

import (
	"lcc1/neoparser"
)

type CodegenResult struct {
	Register string
	TypeInfo neoparser.CompositeType
}
