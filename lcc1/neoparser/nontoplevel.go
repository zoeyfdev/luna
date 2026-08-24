package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
)

func ParseLocal(start int, Tokens []shared.Token, Function *Function) {
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

	ParsePrimary := func() Expression {
		switch peek(0).Type {
		case shared.TokIdent:
			Name := expect(shared.TokIdent)
			IdentObj := Identifier {
				Name: Name,
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

	ParseUnary := func() Expression {
		Expression := ParsePrimary()
		return Expression
	}

	ParseMulDiv := func() Expression {
		Expression := ParseUnary()
		return Expression
	}

	ParseAddSub := func() Expression {
		Expression := ParseMulDiv()
		return Expression	
	}

	ParseExpression := func() {
		LHS := ParseMulDiv()
		switch peek(0).Type {
		case shared.TokEqual:
			// Assignment
			expect(shared.TokEqual)
			RHS := ParseAddSub()

			expect(shared.TokSemi)

			AssignmentObj := Assignment {}
			AssignmentObj.Target = LHS
			AssignmentObj.Value = RHS
			(*Function).Children = append((*Function).Children, AssignmentObj)
		}
	}

	for i < len(Tokens) {
		ParseExpression()
	}
}
