package typecheck

import (
	"lcc1/neoparser"
	"lcc1/shared"
	"lcc1/error"
)

func ReturnTypeName(Type neoparser.CompositeType) string {
	str := ""
	switch Type.Type {
	case neoparser.I8:
		str += "char"
	case neoparser.I16:
		str += "int"
	case neoparser.I32:
		str += "long int"
	default:
		str += "void"
	}

	if Type.PointerLength > 0 {
		str += " "
	}
	for i := 0; i < Type.PointerLength; i++ {
		str += "*"
	}

	if str != "" {
		return str
	}

	return "void"
}

func TypeMediation(T1 TypeCheckReturn, T2 TypeCheckReturn, OpToken shared.Token, Strictness int) TypeCheckReturn {
	// Types:
	// 0: permissive (allows mixing of integers and non-integers)
	// 1: strict (does not allow mixing of integers and non-integers)
	TCR := TypeCheckReturn {
		Token: OpToken,
	}

	var Type1 neoparser.NewType
	var Type2 neoparser.NewType

	Hierarchy := make(map[neoparser.NewType]int)
	Hierarchy[neoparser.I8] = 1
	Hierarchy[neoparser.I16] = 2
	Hierarchy[neoparser.I32] = 3

	if T1.Type.PointerLength > 0 {
		Type1 = neoparser.ReturnUintPtrType()
	} else {
		Type1 = T1.Type.Type
	}

	if T2.Type.PointerLength > 0 {
		Type2 = neoparser.ReturnUintPtrType()
	} else {
		Type2 = T2.Type.Type
	}

	if Strictness == 0 {
		if T1.Type.PointerLength != T2.Type.PointerLength {
			if T1.Type.PointerLength > 0 && T2.Type.PointerLength > 0 {
				error.Error(37, "('" + ReturnTypeName(T1.Type) + "' and '" + ReturnTypeName(T2.Type) + "')", OpToken, T1.TokenSet)
			}
		}
	} else if Strictness == 1 {
		if T1.Type.PointerLength != T2.Type.PointerLength {
			error.Error(37, "('" + ReturnTypeName(T1.Type) + "' and '" + ReturnTypeName(T2.Type) + "')", OpToken, T1.TokenSet)
		}
	}
	
	if Type1 != Type2 {
		if Hierarchy[Type1] > Hierarchy[Type2] {
			TCR.Type.Type = Type1
		} else {
			TCR.Type.Type = Type2
		}
	} else {
		TCR.Type.Type = Type1
	}

	return TCR
}

var CurrentFunction neoparser.Variable

func TypeCheckStatement(Expression neoparser.Expression) {

}

func TypeCheckStatement(Statement neoparser.Statement) neoparser.Statement {
	switch Statement.(type) {
	case neoparser.Assignment:
		Assignment := Statement.(neoparser.Statement)
		return TypeCheckExpression()
	}
}

func TypeCheck(TU *neoparser.AST) {
	for i, Declaration := range (*TU).Declaration {
		if Declaration.Kind == neoparser.FUNCTION {
			CurrentFunction = Declaration
			for j, Statement := range Declaration.Children {
				Statement = TypeCheckStatement(Statement)
			}
		}
	}
}
