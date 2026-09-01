package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
)

var ScopeTicker int = 1
func CreateScope(Parent int) int {
	ID := ScopeTicker
	Scopes = append(Scopes, Scope {
		Parent: Parent,
		ID: ID,
	})
	ScopeTicker++
	return ID
}

func ReturnUintPtrType() NewType {
	switch shared.Bits {
	case 16:
		return I16
	case 32:
		return I32
	}

	return I16
}

func ChildAppend(Slice *[]Statement, Child Statement) {
	*Slice = append(*Slice, Child)
}

func SetRead(Expression Expression, Value bool) Expression {
		switch Expression.(type) {
		case IntLit:
			IntLit := Expression.(IntLit)
			IntLit.IsRead = Value
			return IntLit
		case StringLit:
			StringLit := Expression.(StringLit)
			StringLit.IsRead = Value
			return StringLit
		case Identifier:
			Identifier := Expression.(Identifier)
			Identifier.IsRead = Value
			return Identifier
		case UnaryOperation:
			UnaryOp := Expression.(UnaryOperation)
			UnaryOp.Left = SetRead(UnaryOp.Left, Value)
			return UnaryOp
		case BinaryOperation:
			BinaryOp := Expression.(BinaryOperation)
			BinaryOp.Left = SetRead(BinaryOp.Left, Value)
			BinaryOp.Right = SetRead(BinaryOp.Right, Value)
			return BinaryOp
		case EmptyExpression:
		case StructAccess:
			StructAccess := Expression.(StructAccess)
			StructAccess.IsRead = Value
			return StructAccess
		case IncrementDecrement:
			IncrementDecrement := Expression.(IncrementDecrement)
			return IncrementDecrement
		default:
			error.InternalCompilerError("Unsupported op to SetRead")
		}

		error.InternalCompilerError("No return value for SetRead")
		return IntLit {}
	}
