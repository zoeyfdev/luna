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

func TypeMediation(T1 TypeCheckReturn, T2 TypeCheckReturn, OpToken shared.Token) TypeCheckReturn {
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

	if T1.Type.PointerLength != T2.Type.PointerLength {
		if T1.Type.PointerLength > 0 && T2.Type.PointerLength > 0 {
			error.Error(37, "('PLACEHOLDER' and 'PLACEHOLDER')", OpToken, T1.TokenSet)
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

func TypeCheckLeaf(Leaf neoparser.Leaf) TypeCheckReturn {
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

func TypeCheckUnaryOp(UnaryOp neoparser.UnaryOperation) TypeCheckReturn {
	Left := TypeCheckExpression(UnaryOp.Left)
	switch UnaryOp.Op {
	case shared.TokStar:
		// Dereference
		if Left.OriginalType.PointerLength <= 0 {
			error.Error(26, "", Left.Token, Left.TokenSet)
		}
		Left.Type.PointerLength--
	case shared.TokAmpersand:
		Left.Type.PointerLength++

		if Left.Type.PointerLength > Left.OriginalType.PointerLength + 1 {
			error.Error(27, "'PLACEHOLDER'", Left.Token, Left.TokenSet)
		}
	}

	return Left
}

func TypeCheckBinaryOp(BinaryOp neoparser.BinaryOperation) TypeCheckReturn {
	Left := TypeCheckExpression(BinaryOp.Left)
	Right := TypeCheckExpression(BinaryOp.Right)

	return TypeMediation(Left, Right, BinaryOp.Token)
}

func TypeCheckExpression(Expression neoparser.Expression) TypeCheckReturn {
	switch Expression.(type) {
	case neoparser.Assignment:
		Assignment := Expression.(neoparser.Assignment)
		Target := TypeCheckExpression(Assignment.Target)
		Value := TypeCheckExpression(Assignment.Value)

		TypeMediation(Target, Value, Assignment.Token)
	case neoparser.FunctionCall:
		FunctionCall := Expression.(neoparser.FunctionCall)

		return TypeCheckReturn {
			OriginalType: FunctionCall.AttachedVariable.TypeInfo,
			Type: FunctionCall.AttachedVariable.TypeInfo,
			Token: FunctionCall.Token,
			TokenSet: FunctionCall.TokenSet,
		}	
	case neoparser.BinaryOperation:
		return TypeCheckBinaryOp(Expression.(neoparser.BinaryOperation))
	case neoparser.UnaryOperation:
		return TypeCheckUnaryOp(Expression.(neoparser.UnaryOperation))
	case neoparser.Leaf:
		return TypeCheckLeaf(Expression.(neoparser.Leaf))
	}

	return TypeCheckReturn {}
}

func TypeCheckStatement(Statement neoparser.Statement) {
	switch Statement.(type) {
	case neoparser.Assignment, neoparser.FunctionCall:
		TypeCheckExpression(Statement.(neoparser.Expression))
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
