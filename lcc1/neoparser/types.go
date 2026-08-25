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

type VariableKind int
const (
	FUNCTION VariableKind = iota
	VARIABLE
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

type Variable struct {
	Name string
	TypeInfo CompositeType
	RequiresConstant bool
	Attributes []string
	Internal string
	Scope int
	BasinSize int
	PredefinedValue string
	Parameters []Function_Parameter // hold that thought
	Children []Statement
	Kind VariableKind
}

type Assembly struct {
	String string
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

type FakeStatement struct {
	Value Expression
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

type FunctionCall struct {
	AttachedVariable *Variable
	Children []Expression
}

type Return struct {
	Value Expression
}

// Membership decls
func (_ Variable) _Declaration() {}

func (_ Assembly) _Declaration() {}

func (_ IntLit) _Leaf() {}
func (_ IntLit) _Expression() {}

func (_ Identifier) _Leaf() {}
func (_ Identifier) _Expression() {}

func (_ Assignment) _Statement() {}
func (_ Assignment) _Expression() {}

func (_ BinaryOperation) _Expression() {}

func (_ UnaryOperation) _Expression() {}

func (_ Return) _Statement() {}

func (_ FunctionCall) _Statement() {}
func (_ FunctionCall) _Expression() {}

func (_ FakeStatement) _Statement() {}
