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

type TypeMapKind int
const (
	NORMAL TypeMapKind = iota
	_STRUCT
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
	Children []CompositeType
	MemberName string
	HighName string // TODO: consolidate these
	Offset int
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
	Parameters []Variable // hold that thought
	Children []Statement
	Kind VariableKind
}

type Assembly struct {
	String string
}

// Expressions

type IntLit struct {
	Value string
	Type CompositeType
	Annotated CompositeType
	Token shared.Token
	TokenSet *[]shared.Token
	IsRead bool
}

type StringLit struct {
	Value string
	Type CompositeType
	Annotated CompositeType
	IsRead bool
	Token shared.Token
	TokenSet *[]shared.Token
}

type Identifier struct {
	Name string
	Scope int
	IsRead bool
	Type CompositeType
	Annotated CompositeType
	AttachedVariable Variable
	Token shared.Token
	TokenSet *[]shared.Token
}

type BinaryOperation struct {
	Op shared.TokenType
	Left Expression
	Right Expression
	Type CompositeType
	Token shared.Token
	TokenSet *[]shared.Token
}

type UnaryOperation struct {
	Op shared.TokenType
	Left Expression
	Type CompositeType
}

type Cast struct {
	Value Expression
	Type CompositeType
}

type StructAccess struct {
	Target Expression
	Pointer bool
	Member string
	Type CompositeType
	Token shared.Token
	TokenSet *[]shared.Token
	IsRead bool
}

type EmptyExpression struct {}

type IncrementDecrement struct {
	Target Expression
	Decrement bool
	Type CompositeType
	Token shared.Token
	TokenSet *[]shared.Token
	Post bool
}

// Statements

type ConstAssignStatement struct {
	Value Expression
}

type Assignment struct {
	Target Expression
	Value Expression
	Token shared.Token
	TokenSet *[]shared.Token
}

type FunctionCall struct {
	AttachedVariable Variable
	Children []Expression
	Token shared.Token
	TokenSet *[]shared.Token
	Pushed int
}

type Return struct {
	Value Expression
}

type IfStatement struct {
	Condition Expression
	SuccessChildren []Statement
	ElseChildren []Statement
}

type WhileStatement struct {
	Condition Expression
	Children []Statement
}

type ForStatement struct {
	Condition Expression
	Iterator Expression
	Children []Statement
}

type StatementExpression struct {
	Expression Expression
}

// Membership decls
func (_ Variable) _Declaration() {}

func (_ Assembly) _Declaration() {}
func (_ Assembly) _Statement() {}

func (_ IntLit) _Leaf() {}
func (_ IntLit) _Expression() {}

func (_ Identifier) _Leaf() {}
func (_ Identifier) _Expression() {}

func (_ StringLit) _Leaf() {}
func (_ StringLit) _Expression() {}

func (_ Assignment) _Statement() {}
func (_ Assignment) _Expression() {}

func (_ BinaryOperation) _Expression() {}

func (_ UnaryOperation) _Expression() {}

func (_ Return) _Statement() {}

func (_ FunctionCall) _Statement() {}
func (_ FunctionCall) _Expression() {}

func (_ Cast) _Expression() {}

func (_ ConstAssignStatement) _Statement() {}

func (_ IfStatement) _Statement() {}

func (_ WhileStatement) _Statement() {}

func (_ ForStatement) _Statement() {}

func (_ EmptyExpression) _Expression() {}

func (_ StructAccess) _Expression() {}

func (_ IncrementDecrement) _Expression() {}

func (_ StatementExpression) _Statement() {}

// Other
var TypeMap []CompositeType
