package neoparser

import (
	"lcc1/shared"
)

// Primitives

type NewType int
const (
	VOID NewType = iota
	I8
	I16
	I32
	STRUCT
)

// Interfaces

type Declaration interface {
	_Declaration()
}

type Leaf interface {
	_Leaf()
}

type Expression interface {
	_Expression()
}

type Statement interface {
	_Statement()
}

// Objects

type AST struct {
	Declarations []Declaration
}

type CompositeType struct {
	Type NewType
	PointerLength int
	Size int
	Constant bool
	Static bool
	Signed bool
	Extern bool
}

type Function_Parameter struct {
	Name string
	TypeInfo CompositeType
}

type Function struct {
	Name string
	Parameters []Function_Parameter
	TypeInfo CompositeType
	BasinSize int
	Attributes []string
	Children []Statement
}

type Variable struct {
	Name string
	TypeInfo CompositeType
	RequiresConstant bool
	Attributes []string
	Internal string
}

type IntLit struct {
	Value string
	Scope int
}

type Identifier struct {
	Name string
	Type CompositeType
}

type BinaryOp struct {
	Op shared.Token
	Left Expression
	Right Expression
	Type CompositeType
}

type Assignment struct {
	Target Expression
	Value Expression
}

// List-objects



// Membership decls
func (_ Function) _Declaration() {}
func (_ Variable) _Declaration() {}
func (_ IntLit) _Leaf() {}
func (_ Identifier) _Leaf() {}
func (_ IntLit) _Expression() {}
func (_ Identifier) _Expression() {}
func (_ Assignment) _Expression() {}
func (_ Assignment) _Statement() {}
