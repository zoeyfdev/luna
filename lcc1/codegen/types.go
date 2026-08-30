package codegen

import (
	"lcc1/neoparser"
)

type CodegenResult struct {
	Register string
	TypeInfo neoparser.CompositeType
	ConstantReturn string
	Read bool
	IsRvalue bool
	RValueDerefs int
}
