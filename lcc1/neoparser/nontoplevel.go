package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
	"fmt"
)

func ParseLocal(start int, last int, ScopeID int, Tokens []shared.Token, Children *[]Statement, TU *AST, ExpressionOnly bool, RequireConst bool) {
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

	SearchVariable := func(name string) Variable {
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

				return V
			}
		}

		if ScopeObj.Parent != -1 {
			sid = ScopeObj.Parent
			goto TOP
		}
		
		error.Error(4, "'" + name + "'", peek(-1), &Tokens)
		return Variable {
			Name: "__ZERO",
		}
	}

	var CheckConstant func(Expy Expression, Allow bool)

	CheckConstant = func(Expy Expression, Allow bool) {
		switch Expy.(type) {
		case Identifier:
			Ident := Expy.(Identifier)
			if Allow == false {
				error.Error(35, "", Ident.Token, &Tokens)
			}
		case IntLit:
		case UnaryOperation:
			UnaryOp := Expy.(UnaryOperation)
			switch UnaryOp.Op {
			case shared.TokAmpersand:
				CheckConstant(UnaryOp.Left, true)
			default:
				CheckConstant(UnaryOp.Left, false)
			}
		case BinaryOperation:
			BinaryOp := Expy.(BinaryOperation)
			CheckConstant(BinaryOp.Left, false)
			CheckConstant(BinaryOp.Right, false)
		}
	}

	ParseType := func() (CompositeType, shared.Token) {
		Preset := 0
		PointerLength := 0

		Type := CompositeType {}

		exit := false
		long := false
		short := false

		for i < len(Tokens) {
			switch peek(0).Type {
			case shared.TokQualifier:
				switch peek(0).Value {
				case "short":
					if long == true {
						error.Error(12, "'long' declaration specifier", peek(0), &Tokens)
						break
					}

					if short == false {
						short = true
						Preset = 8
					} else {
						error.Warning(28, "'short' declaration specifier", peek(0), &Tokens)
					}
				case "extern":
					Type.Extern = true
				case "long":
					if short == true {
						error.Error(12, "'short' declaration specifier", peek(0), &Tokens)
						break
					}

					if long == false {
						long = true
						Preset = 32
					} else {
						error.Warning(28, "'short' declaration specifier", peek(0), &Tokens)
					}
				}
				expect(shared.TokQualifier)
			case shared.TokType:
				// qualdone = true
				// typedone = true
				TypeTokString := expect(shared.TokType)
				switch TypeTokString {
				case "int":
					if Preset == 0 {
						switch shared.Bits {
						case 16:
							Type.Type = I16
						case 32:
							Type.Type = I32
						}
					} else {
						switch Preset {
						case 8:
							Type.Type = I8
						case 16:
							Type.Type = I16
						case 32:
							Type.Type = I32
						}
					}
				case "void":
					str := ""
					if short == true {
						str += "short "
					}
					if long == true {
						str += "long "
					}
					str += TypeTokString


					if short == true || long == true {
						error.Error(22, "'" + str + "' is invalid", peek(-1), &Tokens)
					}
					Type.Type = VOID
				case "char":
					str := ""
					if short == true {
						str += "short "
					}
					if long == true {
						str += "long "
					}
					str += TypeTokString


					if short == true || long == true {
						error.Error(22, "'" + str + "' is invalid", peek(-1), &Tokens)
					}

					Type.Type = I8
				}
			case shared.TokStar:
				expect(shared.TokStar)
				PointerLength++
			default:
				exit = true
			}
			if exit == true {
				break
			}
		}

		Type.PointerLength = PointerLength
		
		Size := 0

		if Type.PointerLength > 0 || Type.Type == VOID {
			switch shared.Bits {
			case 16:
				Size = 2
			case 32:
				Size = 4
			}
		} else {
			switch Type.Type {
			case I8:
				Size = 1
			case I16:
				Size = 2
			case I32:
				Size = 4
			}
		}

		Type.Size = Size

		return Type, peek(-1)
	}

	var ParseUnary func(IsRead bool, ForcedType CompositeType) Expression
	var ParseExpression func(IsRead bool, ForcedType CompositeType) Expression

	ParsePrimary := func(IsRead bool, ForcedType CompositeType) Expression {
		switch peek(0).Type {
		case shared.TokLParen:
			switch peek(1).Type {
			case shared.TokQualifier, shared.TokType:
				return ParseUnary(IsRead, ForcedType)
			}
			expect(shared.TokLParen)
			switch peek(0).Type {
			default:
				Expy := ParseExpression(IsRead, ForcedType)
				expect(shared.TokRParen)
				return Expy
			}
			expect(shared.TokRParen)	
		case shared.TokIdent:
			Name := expect(shared.TokIdent)
			ReferenceToken := peek(-1)
			Variable := SearchVariable(Name)

			ActualType := Variable.TypeInfo
			if ForcedType.Type != NONE {
				ActualType = ForcedType
			}

			IdentObj := Identifier {
				IsRead: IsRead,
				Name: Name,
				AttachedVariable: Variable,	
				Type: ActualType,
				Annotated: ActualType,
				Token: ReferenceToken,
				TokenSet: &Tokens,
			}

			switch peek(0).Type {
			case shared.TokLParen:
				// Function call
				expect(shared.TokLParen)
				
				CallObj := FunctionCall {
					Token: peek(-2),
					TokenSet: &Tokens,
				}

				exit := false
				pushed := 0
				for {
					switch peek(0).Type {
					case shared.TokRParen:
						exit = true
						goto DONE
					}

					CallObj.Children = append(CallObj.Children, ParseExpression(true, CompositeType {
						Type: NONE,
					}))
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
			ActualType := CompositeType {
				Type: ReturnUintPtrType(),
			}
			if ForcedType.Type != NONE {
				ActualType = ForcedType
			}

			Number := expect(shared.TokNumber)
			IntObj := IntLit {
				Value: Number,
				Token: peek(-1),
				TokenSet: &Tokens,
				Annotated: ActualType,
				Type: ActualType,
			}
			return IntObj
		case shared.TokString:
			String := expect(shared.TokString)

			ActualType := CompositeType {
				Type: I8,
				PointerLength: 1,
			}

			if ForcedType.Type != NONE {
				ActualType = ForcedType
			}

			StringObj := StringLit {
				Value: String,
				Annotated: ActualType,
				Type: ActualType,
				IsRead: IsRead,
				Token: peek(-1),
				TokenSet: &Tokens,
			}
			return StringObj
		}

		println(peek(0).Value)
		error.InternalCompilerError("no return value")
		return Identifier {
			Name: "__ZERO",
		}
	}

	ParseUnary = func(IsRead bool, ForcedType CompositeType) Expression {
		var Expy Expression = nil

		switch peek(0).Type {
		case shared.TokLParen:
			switch peek(1).Type {
			case shared.TokQualifier, shared.TokType:
				expect(shared.TokLParen)
				Type, _ := ParseType()
				expect(shared.TokRParen)
				return Cast {
					Value: ParseUnary(IsRead, ForcedType),
					Type: Type,
				}
			}
		case shared.TokStar:
			// Dereference
			expect(shared.TokStar)
			Expy = UnaryOperation {
				Op: shared.TokStar,
				Left: ParseUnary(IsRead, ForcedType),
			}	
		case shared.TokAmpersand:
			expect(shared.TokAmpersand)
			Expy = UnaryOperation {
				Op: shared.TokAmpersand,
				Left: ParseUnary(IsRead, ForcedType),
			}
		// TODO: add bang, negative (though L2 at the moment doesn't support negatives)
		}

		if Expy != nil {
			return Expy
		}
		Expression := ParsePrimary(IsRead, ForcedType)
		return Expression
	}

	ParseMulDiv := func(IsRead bool, ForcedType CompositeType) Expression {
		LHS := ParseUnary(IsRead, ForcedType)

		exit := false
		for {
			switch peek(0).Type {
			case shared.TokStar:
				expect(shared.TokStar)
				RHS := ParseUnary(IsRead, ForcedType)
				LHS = BinaryOperation {
					Op: shared.TokStar,
					Left: LHS,
					Right: RHS,
					Token: peek(-1),
				}
			case shared.TokSlash:
				expect(shared.TokSlash)
				RHS := ParseUnary(IsRead, ForcedType)
				LHS = BinaryOperation {
					Op: shared.TokSlash,
					Left: LHS,
					Right: RHS,
					Token: peek(-1),
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

	ParseAddSub := func(IsRead bool, ForcedType CompositeType) Expression {
		LHS := ParseMulDiv(IsRead, ForcedType)

		exit := false
		for {
			switch peek(0).Type {
			case shared.TokPlus:
				expect(shared.TokPlus)
				RHS := ParseMulDiv(IsRead, ForcedType)
				LHS = BinaryOperation {
					Op: shared.TokPlus,
					Left: LHS,
					Right: RHS,
					Token: peek(-1),
				}
			case shared.TokMinus:
				expect(shared.TokMinus)
				RHS := ParseMulDiv(IsRead, ForcedType)
				LHS = BinaryOperation {
					Op: shared.TokMinus,
					Left: LHS,
					Right: RHS,
					Token: peek(-1),
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

	ParseComparator := func(IsRead bool, ForcedType CompositeType) Expression {
		LHS := ParseAddSub(IsRead, ForcedType)
		switch peek(0).Type {
		case shared.TokEquality, shared.TokInequality:
			expect(peek(0).Type)
			Op := peek(-1).Type
			RHS := ParseAddSub(IsRead, ForcedType)	
			return BinaryOperation {
				Left: LHS,
				Right: RHS,
				Op: Op,
			}
		}
	
		return LHS
	}

	ParseExpression = func(IsRead bool, ForcedType CompositeType) Expression {
		LHS := ParseComparator(IsRead, ForcedType)
		switch peek(0).Type {
		case shared.TokEqual:
			// Assignment
			AssignmentObj := Assignment {
				Token: peek(0),
			}

			expect(shared.TokEqual)

			RHS := ParseComparator(true, ForcedType)
			AssignmentObj.Target = LHS
			AssignmentObj.Value = RHS
			return AssignmentObj
		}

		return LHS
	}

	ParseStatement := func() {
		switch peek(0).Type {
		case shared.TokIdent, shared.TokStar, shared.TokString:
			// TODO: change this so it won't crash on non-statement
			*Children = append(*Children, ParseExpression(false, CompositeType {
				Type: NONE,
			}).(Statement))
			expect(shared.TokSemi)
		case shared.TokReturn:
			expect(shared.TokReturn)
			switch peek(0).Type {
			case shared.TokSemi:
			default:
				*Children = append(*Children, Return {
					Value: ParseExpression(true, CompositeType {
						Type: NONE,
					}),
				})
			}
			expect(shared.TokSemi)
		case shared.TokQualifier, shared.TokType:
			slice := []shared.Token {}
			
			exit := false
			for j := i; j < len(Tokens); j++ {
				slice = append(slice, Tokens[j])
				switch Tokens[j].Type {
				case shared.TokSemi:
					i = j + 1
					exit = true
				}
				if exit == true {
					break
				}
			}

			ParseTop(slice, ScopeID, TU, CurrentFunction)
		case shared.TokIf:
			expect(shared.TokIf)
			expect(shared.TokLParen)

			IfObj := IfStatement {}

			IfObj.Condition = ParseExpression(true, CompositeType {
				Type: NONE,
			})
			
			expect(shared.TokRParen)

			expect(shared.TokLCurly)

			slice := []shared.Token {}

			depth := 1
			exit := false
			_Scope := CreateScope(ScopeID)

			for j := i; j < len(Tokens); j++ {
				switch Tokens[j].Type {
				case shared.TokLCurly:
					slice = append(slice, Tokens[j])
					depth++
				case shared.TokRCurly:
					depth--
					if depth <= 0 {
						i = j
						exit = true	
					} else {
						slice = append(slice, Tokens[j])
					}
				default:
					slice = append(slice, Tokens[j])
				}
				if exit == true {
					break
				}
			}

			expect(shared.TokRCurly)

			ParseLocal(0, len(slice) - 1, _Scope, slice, &IfObj.SuccessChildren, TU, false, false)

			*Children = append(*Children, IfObj)
		default:
			error.Error(1, "'" + peek(0).Value + "'", peek(0), &Tokens)
			i++
		}
	}

	if ExpressionOnly == false {
		for i <= last {
			ParseStatement()
		}
	} else {
		Expy := ParseExpression(true, CompositeType {
			Type: NONE,
		})
		CheckConstant(Expy, false)
		*Children = append(*Children, ConstAssignStatement {
			Value: Expy,
		})
	}
}
