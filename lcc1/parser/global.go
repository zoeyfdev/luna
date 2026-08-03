package parser

var TypeMap []TypeMapEntry

var Variables = []Variable_Static {
	{Name: "_r0", Real: "r0", Register: true, Scope: 1, Type: NUMBER32},	
	{Name: "_r1", Real: "r1", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r2", Real: "r2", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r3", Real: "r3", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r4", Real: "r4", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r5", Real: "r5", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r6", Real: "r6", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r7", Real: "r7", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r8", Real: "r8", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r9", Real: "r9", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r10", Real: "r10", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r11", Real: "r11", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_r12", Real: "r12", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e0", Real: "e0", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e1", Real: "e1", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e2", Real: "e2", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e3", Real: "e3", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e4", Real: "e4", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e5", Real: "e5", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e6", Real: "e6", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e7", Real: "e7", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e8", Real: "e8", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e9", Real: "e9", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e10", Real: "e10", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e11", Real: "e11", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e12", Real: "e12", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e13", Real: "e13", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_e14", Real: "e14", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_sp", Real: "sp", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_pc", Real: "pc", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_irv", Real: "irv", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_ir", Real: "ir", Register: true, Scope: 1, Type: NUMBER32},
	{Name: "_b", Real: "b", Register: true, Scope: 1, Type: NUMBER32},	
}

var FunctionDecls = []FunctionDecl {}
var PIE bool

var Scopes = []Scope {
	Scope{ID: 1, Parent: -1},
}
