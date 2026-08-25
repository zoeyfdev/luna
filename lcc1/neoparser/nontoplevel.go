package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
	"fmt"
)

func ParseLocal(start int, ScopeID int, Tokens []shared.Token, Function *Variable, TU *AST, ExpressionOnly bool) {
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
	var ParseExpression func(IsRead bool) Expression

	ParsePrimary := func(IsRead bool) Expression {
		switch peek(0).Type {
		case shared.TokLParen:
			expect(shared.TokLParen)
			switch peek(0).Type {
			default:
				Expy := ParseExpression(IsRead)
				expect(shared.TokRParen)
				return Expy
			}
			expect(shared.TokRParen)
		case shared.TokIdent:
			Name := expect(shared.TokIdent)
			Variable := SearchVariable(Name)

			IdentObj := Identifier {
				IsRead: IsRead,
				Name: Name,
				AttachedVariable: Variable,
			}

			switch peek(0).Type {
			case shared.TokLParen:
				// Function call
				expect(shared.TokLParen)
				
				CallObj := FunctionCall {}

				exit := false
				pushed := 0
				for {
					switch peek(0).Type {
					case shared.TokRParen:
						exit = true
						goto DONE
					}

					CallObj.Children = append(CallObj.Children, ParseExpression(true))
					pushed++

					switch peek(0).Type {
					case shared.TokRParen:
						exit = true
					default:
						expect(shared.TokComma)
					}
					
					DONE:
					if exit == true {
						break
					}
				}

				CallObj.AttachedVariable = Variable

				expect(shared.TokRParen)

				return CallObj
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

		exit := false
		for {
			switch peek(0).Type {
			case shared.TokStar:
				expect(shared.TokStar)
				RHS := ParseUnary(IsRead)
				LHS = BinaryOperation {
					Op: shared.TokStar,
					Left: LHS,
					Right: RHS,
				}
			case shared.TokSlash:
				expect(shared.TokSlash)
				RHS := ParseUnary(IsRead)
				LHS = BinaryOperation {
					Op: shared.TokSlash,
					Left: LHS,
					Right: RHS,
				}
			default:
				exit = true
			}
			if exit == true {
				break
			}
		}
		

		return LHS
	}

	ParseAddSub := func(IsRead bool) Expression {
		LHS := ParseMulDiv(IsRead)

		exit := false
		for {
			switch peek(0).Type {
			case shared.TokPlus:
				expect(shared.TokPlus)
				RHS := ParseMulDiv(IsRead)
				LHS = BinaryOperation {
					Op: shared.TokPlus,
					Left: LHS,
					Right: RHS,
				}
			case shared.TokMinus:
				expect(shared.TokMinus)
				RHS := ParseMulDiv(IsRead)
				LHS = BinaryOperation {
					Op: shared.TokMinus,
					Left: LHS,
					Right: RHS,
				}
			default:
				exit = true
			}
			if exit == true {
				break
			}
		}

		return LHS	
	}

	ParseExpression = func(IsRead bool) Expression {
		LHS := ParseAddSub(IsRead)
		switch peek(0).Type {
		case shared.TokEqual:
			// Assignment
			expect(shared.TokEqual)
			RHS := ParseAddSub(true)

			AssignmentObj := Assignment {}
			AssignmentObj.Target = LHS
			AssignmentObj.Value = RHS
			return AssignmentObj
		}

		return LHS
	}

	ParseStatement := func() {
		switch peek(0).Type {
		case shared.TokIdent, shared.TokStar:
			// TODO: change this so it won't crash on non-statement
			(*Function).Children = append((*Function).Children, ParseExpression(false).(Statement))
			expect(shared.TokSemi)
		case shared.TokReturn:
			expect(shared.TokReturn)
			switch peek(0).Type {
			case shared.TokSemi:
			default:
				(*Function).Children = append((*Function).Children, Return {
					Value: ParseAddSub(true),
				})
			}
			expect(shared.TokSemi)
		default:
			error.Error(1, "'" + peek(0).Value + "'", peek(0), &Tokens)
			i++
		}
	}

	if ExpressionOnly == false {
		for i < len(Tokens) {
			ParseStatement()
		}
	} else {
		(*Function).Children = append((*Function).Children, FakeStatement {
			Value: ParseExpression(true),
		})
	}
}
