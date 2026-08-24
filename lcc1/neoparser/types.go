package neoparser

import (
	"lcc1/shared"
)

// Primitives

type NewType int
const (
	NONE NewType = iota
	VOID
	I8
	I16
	I32
	STRUCT
)

// Interfaces

type Leaf interface {
	_Leaf()
}

type Expression interface {
	_Expression()
}

type Statement interface {
	_Statement()
}

type Declaration interface {
	_Declaration()
}

// Primitive objects

type Scope struct {
	ID int
	Parent int
}

// AST Objects

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
	Scope int
}

type Variable struct {
	Name string
	TypeInfo CompositeType
	RequiresConstant bool
	Attributes []string
	Internal string
	Scope int
	PredefinedValue string
}

type IntLit struct {
	Value string
	Type CompositeType
}

type Identifier struct {
	Name string
	Scope int
	IsRead bool
	AttachedVariable *Variable
}

type BinaryOperation struct {
	Op shared.TokenType
	Left Expression
	Right Expression
	Type CompositeType
}

type UnaryOperation struct {
	Op shared.TokenType
	Left Expression
}

// Statements

type Assignment struct {
	Target Expression
	Value Expression
}

type Return struct {
	Value Expression
}

// Membership decls
func (_ Function) _Declaration() {}
func (_ Variable) _Declaration() {}
func (_ IntLit) _Leaf() {}
func (_ Identifier) _Leaf() {}
func (_ IntLit) _Expression() {}
func (_ Identifier) _Expression() {}
func (_ Assignment) _Statement() {}
func (_ BinaryOperation) _Expression() {}
func (_ UnaryOperation) _Expression() {}
func (_ Return) _Statement() {}
