package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
	"fmt"
)

var ScopeCurrent int
func ParseLocal(start int, last int, ScopeID int, Tokens []shared.Token, Children *[]Statement, TU *AST, ExpressionOnly bool, RequireConst bool) {
	i := start
	ScopeCurrent = ScopeID
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
		sid := ScopeCurrent
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

	ConstructSlice := func() []shared.Token {
		slice := []shared.Token {}

		depth := 1
		exit := false	

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

		return slice
	}

	var SetRead func(Expression Expression, Value bool) Expression
	SetRead = func(Expression Expression, Value bool) Expression {
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
		default:
			error.InternalCompilerError("Unsupported op to SetReadTrue")
		}

		return IntLit {}
	}

	var ParseUnary func(IsRead bool) Expression
	var ParseExpression func(IsRead bool) Expression
	var ParseStatement func(Slice *[]Statement, DefinedScope int)

	ParsePrimary := func(IsRead bool) Expression {
		switch peek(0).Type {
		case shared.TokLParen:
			switch peek(1).Type {
			case shared.TokQualifier, shared.TokType:
				return ParseUnary(IsRead)
			}
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
			ReferenceToken := peek(-1)
			Variable := SearchVariable(Name)

			IdentObj := Identifier {
				IsRead: IsRead,
				Name: Name,
				AttachedVariable: Variable,	
				Type: Variable.TypeInfo,
				Annotated: Variable.TypeInfo,
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
				for {
					switch peek(0).Type {
					case shared.TokRParen:
						exit = true
						goto DONE
					}

					CallObj.Children = append(CallObj.Children, ParseExpression(true))
					CallObj.Pushed++

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
			Type := CompositeType {
				Type: ReturnUintPtrType(),
			}	

			Number := expect(shared.TokNumber)
			IntObj := IntLit {
				Value: Number,
				Token: peek(-1),
				TokenSet: &Tokens,
				Annotated: Type,
				Type: Type,
				IsRead: IsRead,
			}
			return IntObj
		case shared.TokString:
			String := expect(shared.TokString)
			Type := CompositeType {
				Type: I8,
				PointerLength: 1,
			}

			StringObj := StringLit {
				Value: String,
				Annotated: Type,
				Type: Type,
				IsRead: IsRead,
				Token: peek(-1),
				TokenSet: &Tokens,
			}
			return StringObj
		}

		error.InternalCompilerError("no return value")
		return Identifier {
			Name: "__ZERO",
		}
	}

	ParseUnary = func(IsRead bool) Expression {
		var Expy Expression = nil

		switch peek(0).Type {
		case shared.TokLParen:
			switch peek(1).Type {
			case shared.TokQualifier, shared.TokType:
				expect(shared.TokLParen)
				Type, _ := ParseType()
				expect(shared.TokRParen)
				return Cast {
					Value: ParseUnary(IsRead),
					Type: Type,
				}
			}
		case shared.TokStar:
			// Dereference
			expect(shared.TokStar)
			Expy = UnaryOperation {
				Op: shared.TokStar,
				Left: ParseUnary(IsRead),
			}	
		case shared.TokAmpersand:
			expect(shared.TokAmpersand)
			Expy = UnaryOperation {
				Op: shared.TokAmpersand,
				Left: ParseUnary(IsRead),
			}
		// TODO: add bang, negative (though L2 at the moment doesn't support negatives)
		}

		if Expy != nil {
			return Expy
		}
		Expression := ParsePrimary(IsRead)
		return Expression
	}

	ParseMulDiv := func(IsRead bool) Expression {
		LHS := ParseUnary(IsRead)

		exit := false
		for {
			switch peek(0).Type {
			case shared.TokStar, shared.TokSlash:
				expect(peek(0).Type)
				RHS := ParseUnary(IsRead)
				LHS = BinaryOperation {
					Op: peek(-1).Type,
					Left: LHS,
					Right: RHS,
					Token: peek(-1),
					TokenSet: &Tokens,
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
			case shared.TokPlus, shared.TokMinus:
				expect(peek(0).Type)
				Op := peek(-1).Type
				RHS := ParseMulDiv(IsRead)
				LHS = BinaryOperation {
					Op: Op,
					Left: LHS,
					Right: RHS,
					Token: peek(-1),
					TokenSet: &Tokens,
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

	ParseComparator := func(IsRead bool) Expression {
		LHS := ParseAddSub(IsRead)
		switch peek(0).Type {
		case shared.TokEquality, shared.TokInequality, shared.TokLAngle, shared.TokRAngle, shared.TokGEqual, shared.TokLEqual:
			expect(peek(0).Type)
			Op := peek(-1).Type
			RHS := ParseAddSub(IsRead)	
			return BinaryOperation {
				Left: LHS,
				Right: RHS,
				Op: Op,
			}
		}
	
		return LHS
	}

	ParseExpression = func(IsRead bool) Expression {
		LHS := ParseComparator(IsRead)
		switch peek(0).Type {
		case shared.TokEqual:
			// Assignment
			LHS = SetRead(LHS, false)
			AssignmentObj := Assignment {
				Token: peek(0),
				TokenSet: &Tokens,
			}

			expect(shared.TokEqual)

			RHS := ParseComparator(true)
			AssignmentObj.Target = LHS
			AssignmentObj.Value = RHS
			return AssignmentObj
		}

		return LHS
	}

	ParseStatement = func(Slice *[]Statement, DefinedScope int) {
		switch peek(0).Type {
		case shared.TokSemi:
			expect(shared.TokSemi);
		case shared.TokAsm:
			AsmObj := Assembly {}

			expect(shared.TokAsm)
			if peek(0).Value == "volatile" {
				expect(shared.TokQualifier)
			}
			expect(shared.TokLParen)
			String := expect(shared.TokString)
			expect(shared.TokRParen)
			expect(shared.TokSemi)

			AsmObj.String = String

			ChildAppend(Slice, AsmObj)
		case shared.TokReturn:
			expect(shared.TokReturn)
			switch peek(0).Type {
			case shared.TokSemi:
			default:
				ChildAppend(Slice, Return {
					Value: ParseExpression(true),
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

			ParseTop(slice, DefinedScope, TU, CurrentFunction)
		case shared.TokIf:
			IfObj := IfStatement {}
			_Scope := CreateScope(ScopeID)
			_Scope2 := CreateScope(ScopeID)

			expect(shared.TokIf)
			expect(shared.TokLParen)

			IfObj.Condition = ParseExpression(true)
			
			expect(shared.TokRParen)

			switch peek(0).Type {
			case shared.TokLCurly:
				expect(shared.TokLCurly)

				Slice := ConstructSlice()	

				expect(shared.TokRCurly)
				ParseLocal(0, len(Slice) - 1, _Scope, Slice, &IfObj.SuccessChildren, TU, false, false)
			default:
				ParseStatement(&IfObj.SuccessChildren, _Scope)	
			}

			switch peek(0).Type {
			case shared.TokElse:
				expect(shared.TokElse)
				switch peek(0).Type {
				case shared.TokLCurly:
					expect(shared.TokLCurly)

					Slice := ConstructSlice()	

					expect(shared.TokRCurly)
					ParseLocal(0, len(Slice) - 1, _Scope2, Slice, &IfObj.ElseChildren, TU, false, false)
				default:
					ParseStatement(&IfObj.ElseChildren, _Scope2)
				}
			}

			ChildAppend(Slice, IfObj)
		case shared.TokWhile:
			WhileObj := WhileStatement {}
			_Scope := CreateScope(ScopeID)

			expect(shared.TokWhile)
			
			expect(shared.TokLParen)
			
			Condition := ParseExpression(true)
			WhileObj.Condition = Condition

			expect(shared.TokRParen)

			switch peek(0).Type {
			case shared.TokLCurly:
				expect(shared.TokLCurly)
				
				Slice := ConstructSlice()

				expect(shared.TokRCurly)
				ParseLocal(0, len(Slice) - 1, _Scope, Slice, &WhileObj.Children, TU, false, false)
			default:
				ParseStatement(&WhileObj.Children, _Scope)
			}

			ChildAppend(Slice, WhileObj)
		case shared.TokFor:
			ForObj := ForStatement {}
			_Scope := CreateScope(ScopeID)

			expect(shared.TokFor)
				
			expect(shared.TokLParen)
		
			ScopeCurrent = _Scope
			ParseStatement(&CurrentFunction.Children, _Scope) // Parse decl
			ForObj.Condition = ParseExpression(true) // Parse condition
			expect(shared.TokSemi)
			ForObj.Iterator = ParseExpression(true) // Parse iterator
			ScopeCurrent = ScopeID

			expect(shared.TokRParen)

			switch peek(0).Type {
			case shared.TokLCurly:
				expect(shared.TokLCurly)
				
				Slice := ConstructSlice()

				expect(shared.TokRCurly)
				ParseLocal(0, len(Slice) - 1, _Scope, Slice, &ForObj.Children, TU, false, false)
			default:
				ParseStatement(&ForObj.Children, _Scope)
			}	

			ChildAppend(Children, ForObj)
		default:
			// assume expression
			ChildAppend(Slice, ParseExpression(false).(Statement))
			expect(shared.TokSemi)
		}
	}

	if ExpressionOnly == false {
		for i <= last {
			ParseStatement(Children, ScopeID)
		}
	} else {
		Expy := ParseExpression(true)
		CheckConstant(Expy, false)
		ChildAppend(Children, ConstAssignStatement {
			Value: Expy,
		})
	}
}
