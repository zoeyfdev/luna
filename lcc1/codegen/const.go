package codegen

import (
	"lcc1/neoparser"
	"lcc1/error"
	"lcc1/shared"
	"strconv"
	"fmt"
)

func ReturnInt(s string) (int64, bool) {
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

func Const_CodegenLeaf(Leaf neoparser.Leaf, AllowIdent bool) string {
	switch Leaf.(type) {
	case neoparser.IntLit:
		IntLit := Leaf.(neoparser.IntLit)
		return IntLit.Value
	case neoparser.StringLit:
		StringLit := Leaf.(neoparser.StringLit)
		VT := VarTicker
		VarTicker++
		WritePre(fmt.Sprintf("const_str_%d:", VT), false)
		WritePre(".asciz \"" + StringLit.Value + "\"\n", true)
		return fmt.Sprintf("const_str_%d", VT)
	case neoparser.Identifier:
		Identifier := Leaf.(neoparser.Identifier)

		if AllowIdent == false {
			error.InternalCompilerError("illegal ident value in Const_CodegenLeaf")
		} else {
			return Identifier.AttachedVariable.Internal
		}
	}

	return ""
}

func Const_CodegenUnaryOp(UnaryOp neoparser.UnaryOperation, AllowIdent bool) string {	
	switch UnaryOp.Op {
	case shared.TokAmpersand:
		return Const_CodegenExpression(UnaryOp.Left, true)	
	default:
		error.InternalCompilerError("invalid unary op type!")
	}
	return Const_CodegenExpression(UnaryOp.Left, AllowIdent)
}

func Const_CodegenBinaryOp(BinaryOp neoparser.BinaryOperation, AllowIdent bool) string {
	LHS, ok := ReturnInt(Const_CodegenExpression(BinaryOp.Left, AllowIdent))
	RHS, ok2 := ReturnInt(Const_CodegenExpression(BinaryOp.Right, AllowIdent))

	if ok == false || ok2 == false {
		error.InternalCompilerError("Fail return number for expression")
	}

	switch BinaryOp.Op {
	case shared.TokPlus:
		return fmt.Sprintf("%d", LHS + RHS)
	case shared.TokMinus:
		return fmt.Sprintf("%d", LHS - RHS)	
	case shared.TokStar:
		return fmt.Sprintf("%d", LHS * RHS)
	case shared.TokSlash:
		return fmt.Sprintf("%d", LHS / RHS)
	default:
		error.InternalCompilerError("invalid binary op type!")
	}

	return ""
}

func Const_CodegenExpression(Expression neoparser.Expression, AllowIdent bool) string {
	switch Expression.(type) {
	case neoparser.BinaryOperation:
		return Const_CodegenBinaryOp(Expression.(neoparser.BinaryOperation), AllowIdent)
	case neoparser.UnaryOperation:
		return Const_CodegenUnaryOp(Expression.(neoparser.UnaryOperation), AllowIdent)
	case neoparser.IntLit, neoparser.Identifier, neoparser.StringLit:
		return Const_CodegenLeaf(Expression.(neoparser.Leaf), AllowIdent)
	}
	return ""
}
