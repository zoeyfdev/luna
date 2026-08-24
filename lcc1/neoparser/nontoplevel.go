package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
	"fmt"
)

func ParseLocal(start int, ScopeID int, Tokens []shared.Token, Function *Function, TU *AST) {
	i := start
	expect := func(toktype shared.TokenType) string {
		var value string
		if i >= len(Tokens) {
			if toktype != shared.TokSemi {
				error.Error(1, "'<EOF>'", Tokens[i - 1], &Tokens)
			} else {
				error.Error(18, "", Tokens[i - 1], &Tokens)
			}
			return ""
		}
		if Tokens[i].Type == toktype {
			value = Tokens[i].Value
		} else {
			if toktype != shared.TokSemi {
				error.Error(1, "'" + Tokens[i].Value + "'", Tokens[i], &Tokens)
			} else {
				error.Error(18, "", Tokens[i - 1], &Tokens)
			}
		}
		i++
		return value
	}	
	peek := func(lookahead int) shared.Token {
		if i + lookahead < len(Tokens) {
			return Tokens[i + lookahead]
		}
		return shared.Token{Type: shared.TokEOF, Value: ""}
	
	}

	SearchVariable := func(name string) *Variable {
		fmt.Println("scope", ScopeID)
		sid := ScopeID
	TOP:
		ScopeObj := Scope {}
		
		Found := false
		for _, S := range Scopes {
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
			case Variable:
				V := Declaration.(Variable)
				if V.Scope != sid || V.Name != name { continue }

				return &V
			}
		}

		fmt.Println(ScopeObj.ID, ScopeObj.Parent)
		if ScopeObj.Parent != -1 {
			sid = ScopeObj.Parent
			goto TOP
		}
		
		error.Error(4, "'" + name + "'", peek(-1), &Tokens)
		return &Variable {
			Name: "__ZERO",
		}
	}

	var ParseUnary func(IsRead bool) Expression

	ParsePrimary := func(IsRead bool) Expression {
		switch peek(0).Type {
		case shared.TokIdent:
			Name := expect(shared.TokIdent)

			IdentObj := Identifier {
				IsRead: IsRead,
				Name: Name,
				AttachedVariable: SearchVariable(Name),
			}
			return IdentObj
		case shared.TokNumber:
			Number := expect(shared.TokNumber)
			IntObj := IntLit {
				Value: Number,
			}
			return IntObj
		}

		error.InternalCompilerError("no return value")
		return Identifier {
			Name: "__ZERO",
		}
	}

	ParseUnary = func(IsRead bool) Expression {
		switch peek(0).Type {
		case shared.TokStar:
			// Dereference
			expect(shared.TokStar)
			return UnaryOperation {
				Op: shared.TokStar,
				Left: ParseUnary(IsRead),
			}
		// TODO: add bang, negative (though L2 at the moment doesn't support negatives)
		}

		Expression := ParsePrimary(IsRead)
		return Expression
	}

	ParseMulDiv := func(IsRead bool) Expression {
		LHS := ParseUnary(IsRead)
		switch peek(0).Type {
		case shared.TokStar:
			expect(shared.TokStar)
			RHS := ParseUnary(IsRead)
			return BinaryOperation {
				Op: shared.TokStar,
				Left: LHS,
				Right: RHS,
			}
		case shared.TokSlash:
			expect(shared.TokSlash)
			RHS := ParseUnary(IsRead)
			return BinaryOperation {
				Op: shared.TokSlash,
				Left: LHS,
				Right: RHS,
			}	
		}

		return LHS
	}

	ParseAddSub := func(IsRead bool) Expression {
		LHS := ParseMulDiv(IsRead)
		switch peek(0).Type {
		case shared.TokPlus:
			expect(shared.TokPlus)
			RHS := ParseMulDiv(IsRead)
			return BinaryOperation {
				Op: shared.TokPlus,
				Left: LHS,
				Right: RHS,
			}
		case shared.TokMinus:
			expect(shared.TokMinus)
			RHS := ParseMulDiv(IsRead)
			return BinaryOperation {
				Op: shared.TokPlus,
				Left: LHS,
				Right: RHS,
			}	
		}
		return LHS	
	}

	ParseExpression := func(IsRead bool) {
		LHS := ParseMulDiv(IsRead)
		switch peek(0).Type {
		case shared.TokEqual:
			// Assignment
			expect(shared.TokEqual)
			RHS := ParseAddSub(true)
			expect(shared.TokSemi)

			AssignmentObj := Assignment {}
			AssignmentObj.Target = LHS
			AssignmentObj.Value = RHS
			(*Function).Children = append((*Function).Children, AssignmentObj)
		}
	}

	ParseStatement := func() {
		switch peek(0).Type {
		case shared.TokIdent, shared.TokStar:
			ParseExpression(false)
		case shared.TokReturn:
			expect(shared.TokReturn)
			(*Function).Children = append((*Function).Children, Return {
				Value: ParseAddSub(true),
			})
			expect(shared.TokSemi)
		}
	}

	for i < len(Tokens) {
		ParseStatement()
	}
}
