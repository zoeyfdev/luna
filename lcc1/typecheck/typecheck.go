package typecheck

import (
	"lcc1/neoparser"
	"lcc1/shared"
	"lcc1/error"
	"reflect"
	"fmt"
)

var CurrentFunction neoparser.Variable

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
	case neoparser.VOID:
		str += "void"
	default:
		str += Type.HighName
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
	// TODO: different error messages for different scenarios
	TCR := TypeCheckReturn {
		Token: OpToken,
	}

	IsInt := func(NT neoparser.NewType) bool {
		switch NT {
		case neoparser.I8, neoparser.I16, neoparser.I32:
			return true
		}
		return false
	}

	var Type1 neoparser.NewType
	var Type2 neoparser.NewType

	Hierarchy := make(map[neoparser.NewType]int)
	Hierarchy[neoparser.I8] = 1
	Hierarchy[neoparser.I16] = 2
	Hierarchy[neoparser.I32] = 3

	if T1.Type.PointerLength > 0 && T2.Type.PointerLength > 0 {
		if T1.Type.Type != T2.Type.Type && !(T1.Type.Type == neoparser.VOID || T2.Type.Type == neoparser.VOID) {
			error.Error(37, "('" + ReturnTypeName(T1.Type) + "' and '" + ReturnTypeName(T2.Type) + "')", OpToken, T1.TokenSet)
		} 
	}

	if (!IsInt(T1.Type.Type) || !IsInt(T2.Type.Type)) && (T1.Type.PointerLength <= 0 || T2.Type.PointerLength <= 0) {
		if (T1.Type.Type != T2.Type.Type) || (T1.Type.HighName != T2.Type.HighName) {
			error.Error(37, "('" + ReturnTypeName(T1.Type) + "' and '" + ReturnTypeName(T2.Type) + "')", OpToken, T1.TokenSet)	
		}
	}

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
			TCR.Type.PointerLength = max(T1.Type.PointerLength, T2.Type.PointerLength)
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

func TypeSweep(Expression neoparser.Expression, Type neoparser.CompositeType) neoparser.Expression {
	switch Expression.(type) {
	case neoparser.IntLit:
		IntLit := Expression.(neoparser.IntLit)
		IntLit.Annotated = Type
		return IntLit
	case neoparser.StringLit:
		StringLit := Expression.(neoparser.StringLit)
		StringLit.Annotated = Type
		return StringLit
	case neoparser.Identifier:
		Identifier := Expression.(neoparser.Identifier)
		Identifier.Annotated = Type
		return Identifier
	case neoparser.UnaryOperation:
		UnaryOp := Expression.(neoparser.UnaryOperation)
		UnaryOp.Left = TypeSweep(UnaryOp.Left, Type)
		return UnaryOp
	case neoparser.BinaryOperation:
		BinaryOp := Expression.(neoparser.BinaryOperation)
		BinaryOp.Left = TypeSweep(BinaryOp.Left, Type)
		BinaryOp.Right = TypeSweep(BinaryOp.Right, Type)
		return BinaryOp
	case neoparser.FunctionCall:
		FunctionCall := Expression.(neoparser.FunctionCall)
		FunctionCall.AttachedVariable.TypeInfo = Type
		return FunctionCall
	case neoparser.StructAccess:
		StructAccess := Expression.(neoparser.StructAccess)
		StructAccess.Type = Type
		return StructAccess
	// ^ added this back
	case neoparser.Subscript:
		Subscript := Expression.(neoparser.Subscript)
		Subscript.Type = Type
		return Subscript
	}
	
	error.InternalCompilerError("no return value for " + reflect.TypeOf(Expression).String())
	return neoparser.IntLit{}
}

func TypeCheckLeaf(Leaf neoparser.Leaf, Strictness int) TypeCheckReturn {
	switch Leaf.(type) {
	case neoparser.IntLit:
		IntLit := Leaf.(neoparser.IntLit)
		return TypeCheckReturn {
			Type: IntLit.Annotated,
			Token: IntLit.Token,
			TokenSet: IntLit.TokenSet,
			Expression: IntLit,
			RValue: true,
		}
	case neoparser.StringLit:
		StringLit := Leaf.(neoparser.StringLit)
		return TypeCheckReturn {
			Type: StringLit.Annotated,
			Token: StringLit.Token,
			TokenSet: StringLit.TokenSet,
			Expression: StringLit,
		}
	case neoparser.Identifier:
		Identifier := Leaf.(neoparser.Identifier)
		return TypeCheckReturn {
			Type: Identifier.Annotated,
			Token: Identifier.Token,
			TokenSet: Identifier.TokenSet,
			Expression: Identifier,
		}
	}

	return TypeCheckReturn {}
}

func TypeCheckUnaryOp(UnaryOp neoparser.UnaryOperation, Strictness int) TypeCheckReturn {
	Left := TypeCheckExpression(UnaryOp.Left, Strictness)
	
	if Left.Expression != nil { UnaryOp.Left = Left.Expression }

	switch UnaryOp.Op {
	case shared.TokStar:
		// Dereference
		if Left.Type.PointerLength <= 0 {
			error.Error(26, "('" + ReturnTypeName(Left.Type) + "' invalid)", UnaryOp.Token, UnaryOp.TokenSet)
		}
		Left.Type.PointerLength--
		Left.RValue = false
	case shared.TokAmpersand:
		Left.AddrCount++

		if Left.AddrCount > Left.Type.PointerLength + 1 {
			error.Error(27, "'" + ReturnTypeName(Left.Type) + "'", UnaryOp.Token, UnaryOp.TokenSet)
		}
		Left.Type.PointerLength++
		Left.RValue = true
	}

	UnaryOp.Type = Left.Type
	Left.Expression = UnaryOp

	return Left
}

func TypeCheckBinaryOp(BinaryOp neoparser.BinaryOperation, Strictness int) TypeCheckReturn {
	Left := TypeCheckExpression(BinaryOp.Left, Strictness)
	Right := TypeCheckExpression(BinaryOp.Right, Strictness)

	BinaryOp.Left = Left.Expression
	BinaryOp.Right = Right.Expression

	TCR := TypeMediation(Left, Right, BinaryOp.Token, Strictness)
	TCR.Expression = BinaryOp
	TCR.RValue = true
	TCR.Token = BinaryOp.Token
	TCR.TokenSet = BinaryOp.TokenSet

	return TCR
}

func TypeCheckExpression(Expression neoparser.Expression, Strictness int) TypeCheckReturn {
	switch Expression.(type) {
	case neoparser.Assignment:
		Assignment := Expression.(neoparser.Assignment)
		Target := TypeCheckExpression(Assignment.Target, Strictness)
		Value := TypeCheckExpression(Assignment.Value, Strictness)

		if Target.RValue == true {
			error.Error(45, "", Assignment.Token, Assignment.TokenSet)
		}

		if Target.Expression != nil { Assignment.Target = Target.Expression }
		if Value.Expression != nil { Assignment.Value = Value.Expression }

		ReturnStmt := TypeMediation(Target, Value, Assignment.Token, 1)
		ReturnStmt.Expression = Assignment
		return ReturnStmt
	case neoparser.FunctionCall:
		FunctionCall := Expression.(neoparser.FunctionCall)

		Expected := len(FunctionCall.AttachedVariable.Parameters)
		if FunctionCall.Pushed < Expected {
			error.Error(20, fmt.Sprintf("expected %d, have %d", Expected, FunctionCall.Pushed), FunctionCall.Token, FunctionCall.TokenSet)
		} else if FunctionCall.Pushed > Expected {
			error.Error(21, fmt.Sprintf("expected %d, have %d", Expected, FunctionCall.Pushed), FunctionCall.Token, FunctionCall.TokenSet)
		}

		for i, Expy := range FunctionCall.Children {
			if i >= len(FunctionCall.AttachedVariable.Parameters) { break }

			Token, _ := neoparser.ReturnTokenPair(Expy)
			TypeMediation(TypeCheckExpression(Expy, 1), TypeCheckReturn {
				Type: FunctionCall.AttachedVariable.Parameters[i].TypeInfo,
			}, Token, 1)
		}
		
		return TypeCheckReturn {
			Expression: FunctionCall,
			Type: FunctionCall.AttachedVariable.TypeInfo,
			Token: FunctionCall.Token,
			TokenSet: FunctionCall.TokenSet,
			RValue: true,
		} 
	case neoparser.BinaryOperation:
		return TypeCheckBinaryOp(Expression.(neoparser.BinaryOperation), Strictness)
	case neoparser.UnaryOperation:
		return TypeCheckUnaryOp(Expression.(neoparser.UnaryOperation), Strictness)
	case neoparser.Leaf:
		return TypeCheckLeaf(Expression.(neoparser.Leaf), Strictness)
	case neoparser.Cast:
		Cast := Expression.(neoparser.Cast)
		Expression = TypeSweep(Cast.Value, Cast.Type)
		return TypeCheckReturn {
			Expression: Expression,
			Type: Cast.Type,
			Token: Cast.Token,
			TokenSet: Cast.TokenSet,
		}
	case neoparser.StructAccess:
		StructAccess := Expression.(neoparser.StructAccess)
		return TypeCheckReturn {
			Expression: StructAccess,
			Type: StructAccess.Type,
			Token: StructAccess.Token,
			TokenSet: StructAccess.TokenSet,
			RValue: false,
		}
	case neoparser.IncrementDecrement:
		IncrementDecrement := Expression.(neoparser.IncrementDecrement)
		Target := TypeCheckExpression(IncrementDecrement.Target, Strictness)

		if Target.RValue == true {
			error.Error(45, "", IncrementDecrement.Token, IncrementDecrement.TokenSet)
		}

		return TypeCheckReturn {
			Expression: IncrementDecrement,
			Type: Target.Type,
			Token: IncrementDecrement.Token,
			TokenSet: IncrementDecrement.TokenSet,
			RValue: true,	
		}
	case neoparser.Subscript:
		Subscript := Expression.(neoparser.Subscript)
		Target := TypeCheckExpression(Subscript.Target, Strictness)


		// TODO: change to allow 0[arr] to work
		// TODO: add arrays to this
		if Target.Type.PointerLength <= 0 {
			error.Error(40, "", Subscript.Token, Subscript.TokenSet)
		}
		
		return TypeCheckReturn {
			Expression: Subscript,
			Type: Subscript.Type,
			Token: Subscript.Token,
			TokenSet: Subscript.TokenSet,
			RValue: false,
		}
	}

	return TypeCheckReturn {}
}

func TypeCheckStatement(Statement neoparser.Statement) {
	switch Statement.(type) {
	case neoparser.Assignment, neoparser.FunctionCall:
		TypeCheckExpression(Statement.(neoparser.Expression), 0)
	case neoparser.StatementExpression:
		StatementExpression := Statement.(neoparser.StatementExpression)
		TypeCheckExpression(StatementExpression.Expression, 0)
	case neoparser.ForStatement:
		ForStatement := Statement.(neoparser.ForStatement)
		for _, Statement := range ForStatement.Initializer {
			TypeCheckStatement(Statement)
		}
		TypeCheckExpression(ForStatement.Condition, 0)
		TypeCheckExpression(ForStatement.Iterator, 0)
		for _, Statement := range ForStatement.Children {
			TypeCheckStatement(Statement)
		}
	case neoparser.WhileStatement:
		WhileStatement := Statement.(neoparser.WhileStatement)
		TypeCheckExpression(WhileStatement.Condition, 0)
		for _, Statement := range WhileStatement.Children {
			TypeCheckStatement(Statement)
		}
	case neoparser.IfStatement:
		IfStatement := Statement.(neoparser.IfStatement)
		TypeCheckExpression(IfStatement.Condition, 0)
		for _, Statement := range IfStatement.SuccessChildren {
			TypeCheckStatement(Statement)
		}
		for _, Statement := range IfStatement.ElseChildren {
			TypeCheckStatement(Statement)
		}
	case neoparser.Return:
		ReturnStatement := Statement.(neoparser.Return)
	
		if ReturnStatement.Value != nil && (CurrentFunction.TypeInfo.Type == neoparser.VOID && CurrentFunction.TypeInfo.PointerLength <= 0) {
			// void func returning a value
			error.Error(44, "'" + CurrentFunction.Name + "' should not return a value", ReturnStatement.Token, ReturnStatement.TokenSet)
		} else if ReturnStatement.Value == nil && !(CurrentFunction.TypeInfo.Type == neoparser.VOID && CurrentFunction.TypeInfo.PointerLength <= 0) {
			// non void func not returning a value
			// note it doesn't apply if a non-void function doesn't include a return at all
			error.Error(56, "", ReturnStatement.Token, ReturnStatement.TokenSet)
		} else {
			if !(CurrentFunction.TypeInfo.Type == neoparser.VOID && CurrentFunction.TypeInfo.PointerLength <= 0) {
				TypeMediation(TypeCheckExpression(ReturnStatement.Value, 1), TypeCheckReturn {
					Type: CurrentFunction.TypeInfo,
				}, ReturnStatement.Token, 1)
			}
		}	
	}	
}

func TypeCheck(TranslationUnit *neoparser.AST) {
	for _, Declaration := range (*TranslationUnit).Declarations {
		switch Declaration.(type) {
		case neoparser.Variable:
			Var := Declaration.(neoparser.Variable)
			switch Var.Kind {
			case neoparser.FUNCTION:
				CurrentFunction = Var
				// We subtract by one to find the actual value, not the value it decays to when raw.
				CurrentFunction.TypeInfo.PointerLength-- 
				for _, Child := range Var.Children {
					TypeCheckStatement(Child)
				}
			}
		}
	} 
}
