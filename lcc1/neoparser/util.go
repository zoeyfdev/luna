package neoparser

import (
	"lcc1/shared"
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

