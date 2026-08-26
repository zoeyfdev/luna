package typecheck

import (
	"lcc1/neoparser"
	"lcc1/shared"
	"lcc1/error"
)

func ReturnUintPtrType() neoparser.NewType {
	switch shared.Bits {
	case 16:
		return neoparser.I16
	case 32:
		return neoparser.I32
	}

	return neoparser.I16
}

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
		Type1 = ReturnUintPtrType()
	} else {
		Type1 = T1.Type.Type
	}

	if T2.Type.PointerLength > 0 {
		Type2 = ReturnUintPtrType()
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

func TypeCheckLeaf(Leaf neoparser.Leaf, Strictness int) TypeCheckReturn {
	switch Leaf.(type) {
	case neoparser.IntLit:
		IntLit := Leaf.(neoparser.IntLit)
		return TypeCheckReturn {
			OriginalType: IntLit.Type,
			Type: IntLit.Type,
			Token: IntLit.Token,
			TokenSet: IntLit.TokenSet,
		}
	case neoparser.StringLit:
		StringLit := Leaf.(neoparser.StringLit)
		return TypeCheckReturn {
			OriginalType: StringLit.Type,
			Type: StringLit.Type,
			Token: StringLit.Token,
			TokenSet: StringLit.TokenSet,
		}
	case neoparser.Identifier:
		Identifier := Leaf.(neoparser.Identifier)
		return TypeCheckReturn {
			OriginalType: Identifier.Type,
			Type: Identifier.Type,
			Token: Identifier.Token,
			TokenSet: Identifier.TokenSet,
		}
	}

	return TypeCheckReturn {}
}

func TypeCheckUnaryOp(UnaryOp neoparser.UnaryOperation, Strictness int) TypeCheckReturn {
	Left := TypeCheckExpression(UnaryOp.Left, Strictness)
	switch UnaryOp.Op {
	case shared.TokStar:
		// Dereference
		if Left.Type.PointerLength <= 0 {
			error.Error(26, "('" + ReturnTypeName(Left.Type) + "' invalid)", Left.Token, Left.TokenSet)
		}
		Left.Type.PointerLength--
	case shared.TokAmpersand:
		Left.Type.PointerLength++

		if Left.Type.PointerLength > Left.OriginalType.PointerLength + 1 {
			error.Error(27, "'" + ReturnTypeName(Left.Type) + "'", Left.Token, Left.TokenSet)
		}
	}

	return Left
}

func TypeCheckBinaryOp(BinaryOp neoparser.BinaryOperation, Strictness int) TypeCheckReturn {
	Left := TypeCheckExpression(BinaryOp.Left, Strictness)
	Right := TypeCheckExpression(BinaryOp.Right, Strictness)

	return TypeMediation(Left, Right, BinaryOp.Token, Strictness)
}

func TypeCheckExpression(Expression neoparser.Expression, Strictness int) TypeCheckReturn {
	switch Expression.(type) {
	case neoparser.Assignment:
		Assignment := Expression.(neoparser.Assignment)
		Target := TypeCheckExpression(Assignment.Target, Strictness)
		Value := TypeCheckExpression(Assignment.Value, Strictness)

		TypeMediation(Target, Value, Assignment.Token, Strictness)
	case neoparser.FunctionCall:
		FunctionCall := Expression.(neoparser.FunctionCall)

		return TypeCheckReturn {
			OriginalType: FunctionCall.AttachedVariable.TypeInfo,
			Type: FunctionCall.AttachedVariable.TypeInfo,
			Token: FunctionCall.Token,
			TokenSet: FunctionCall.TokenSet,
		}	
	case neoparser.BinaryOperation:
		return TypeCheckBinaryOp(Expression.(neoparser.BinaryOperation), Strictness)
	case neoparser.UnaryOperation:
		return TypeCheckUnaryOp(Expression.(neoparser.UnaryOperation), Strictness)
	case neoparser.Leaf:
		return TypeCheckLeaf(Expression.(neoparser.Leaf), Strictness)
	}

	return TypeCheckReturn {}
}

func TypeCheckStatement(Statement neoparser.Statement) {
	switch Statement.(type) {
	case neoparser.Assignment, neoparser.FunctionCall:
		TypeCheckExpression(Statement.(neoparser.Expression), 0)
	}
}

func TypeCheck(TranslationUnit *neoparser.AST) {
	for _, Declaration := range (*TranslationUnit).Declarations {
		Var := Declaration.(neoparser.Variable)
		switch Var.Kind {
		case neoparser.FUNCTION:
			for _, Statement := range Var.Children {
				TypeCheckStatement(Statement)
			}	
		}
	} 
}
