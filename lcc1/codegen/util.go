package codegen

import (
	"lcc1/neoparser"
	"lcc1/error"
	"lcc1/shared"
	"fmt"
)

var VarTicker int

func SearchVariable(name string, ScopeID int, TU *neoparser.AST) neoparser.Variable {
	sid := ScopeID
TOP:
	ScopeObj := neoparser.Scope {}
	
	Found := false
	for _, S := range neoparser.Scopes {
		if S.ID == sid {
			ScopeObj = S
			Found = true
			break
		}
	}

	if Found == false {
		error.InternalCompilerError("scope '" + fmt.Sprintf("%d", ScopeID) + "' not found")
	}

	for _, Declaration := range (*TU).Declarations {
		switch Declaration.(type) {
		case neoparser.Variable:
			V := Declaration.(neoparser.Variable)
			if V.Scope != sid || V.Name != name { continue }

			return V
		}
	}

	if ScopeObj.Parent != -1 {
		sid = ScopeObj.Parent
		goto TOP
	}
	
	return neoparser.Variable {
		Name: "__ZERO",
	}
}

func ReturnUintPtrType() neoparser.NewType {
	switch shared.Bits {
	case 16:
		return neoparser.I16
	case 32:
		return neoparser.I32
	}

	return neoparser.I16
}
