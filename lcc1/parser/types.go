package parser

import (
	"lcc1/shared"
)

var level int = 0
var L1_ALLOW_NONCONST bool = false

var Code1 string = ""
var Code2 string = ""

var IDCounter = 1

const (
	NUMBER8 int = iota // unsigned short short int
	NUMBER16		   // unsigned short int / unsigned int
	NUMBER32           // unsigned long int
	STRING             // unsigned char
	STRUCT              // structs
	NULL               // void / void*
)

type Variable_Static struct {
	Name string
	Type int
	Type2 int
	Value any
	Pointer bool
	Real string
	Scope int
	Const bool
	Extern bool
	ArgNum int
	Register bool
	PointerLength int
	BasinSize uint32
	HasBasin bool
	ArgumentTypeManifest []ArgumentTypeManifestEntry
	StructTotalSize int
	StructMemberList []StructMemberListEntry
	PointingToStruct string 
}


type StructMemberListEntry struct {
	Name string
	Offset int
	RequiredType ArgumentTypeManifestEntry
	PointingToStruct string
}

type ArgumentTypeManifestEntry struct {
	Type int
	Type2 int
	Pointer bool
	PointerLength int
}

type FunctionDecl struct {
	Name string
	Token shared.Token
	Set []shared.Token	
}

type Scope struct {
	ID int
	Parent int
}

type UnpackOrder struct {
	Register string
	Label string
	Type int
	Pointer bool
	PointerLength int
}

type TypeMapEntry struct {
	Name string
	ReferralType int
	Pointer bool
	PointerLength int
	Const bool
	EmbeddedStruct Variable_Static
}

type NewType int

const (
	I8 = iota
	I16
	I32
	STRUCT
)


// AST stuff

type Node interface {
	_Node() int
}

type Declaration interface {
	_Declaration() int
}

// Types

type Function_Parameter struct {
	Name string
	Type NewType
}

type Function struct {
	Name string
	Parameters []Function_Parameter
	Children []Node
}

// Node-types

type SymOp struct {
	
}


// Membership decls
func (_ SymOp) _Node() int { return 0 }
func (_ Function) _Declaration() int { return 0 }
