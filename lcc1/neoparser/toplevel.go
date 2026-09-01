package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
	"fmt"
)

var IDCounter = 1

func Parse(Tokens []shared.Token) *AST {
	TranslationUnit := AST {}

	ParseTop(Tokens, 0, &TranslationUnit, &Variable {}, &[]Statement {})

	return &TranslationUnit
}

func ParseTop(Tokens []shared.Token, Scope int, TU *AST, EnclosingFunction *Variable, VInitChildrenSlice *[]Statement) { // TODO: add dynamic scoping for DECL	
	i := 0
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

	ParseType := func() (CompositeType, shared.Token) {
		Preset := 0
		PointerLength := 0

		Type := CompositeType {}

		exit := false

		// qualdone := false
		// typedone := false

		long := false
		short := false

		for i < len(Tokens) {
			switch peek(0).Type {
			case shared.TokIdent:
				// Custom type
				Name := expect(shared.TokIdent)

				for _, Type := range TypeMap {
					if Type.HighName == Name {
						return Type, peek(-1)
					}
				}

				error.Error(46, "'" + Name + "'", peek(-1), &Tokens)
				exit = true
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
				case "static":
					Type.Static = true
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
				switch peek(0).Type {
				case shared.TokStar:
				default:
					exit = true
				}
			case shared.TokStar:
				expect(shared.TokStar)
				PointerLength++

				switch peek(0).Type {
				case shared.TokStar:
				default:
					exit = true
				}
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

	GenerateLocalIVN := func(EnclosingFunction *Variable, TypeInfo CompositeType) string {
		CurrentOffset := (*EnclosingFunction).BasinSize
		if TypeInfo.PointerLength > 0 || TypeInfo.Type == VOID {
			switch shared.Bits {
			case 16:
				(*EnclosingFunction).BasinSize += 2	
			case 32:
				(*EnclosingFunction).BasinSize += 4
			}
		} else {
			switch TypeInfo.Type {
			case I8:
				(*EnclosingFunction).BasinSize += 1
			case I16:
				(*EnclosingFunction).BasinSize += 2
			case I32:
				(*EnclosingFunction).BasinSize += 4
			case STRUCT:
				(*EnclosingFunction).BasinSize += TypeInfo.Size
			}
		}

		return fmt.Sprintf("fp + %d", CurrentOffset)
	}

	DeclAppend := func(Decl Declaration) {
		(*TU).Declarations = append((*TU).Declarations, Decl)
	} 

	for i < len(Tokens) {
		switch peek(0).Type {
		case shared.TokQualifier, shared.TokType, shared.TokIdent:
			TypeInformation, TypeTok := ParseType()
			NameLocation := i
			Name := expect(shared.TokIdent)

			switch peek(0).Type {
			case shared.TokLParen:
				FObj := Variable {}
				FObj.Name = Name
				FObj.Internal = Name
				FObj.TypeInfo = TypeInformation
				FObj.Kind = FUNCTION
				Scope := CreateScope(0)

				expect(shared.TokLParen)

				exit := false
				for {
					switch peek(0).Type {
					default:
						if peek(0).Type == shared.TokEllipsis {
							error.UnimplementedMessage("ellipsis operators are currently not supported.")
						}
						Type, _ := ParseType() // TODO: integrated checking for void type
						Name := expect(shared.TokIdent)
						if peek(0).Type != shared.TokRParen {
							expect(shared.TokComma)
						}
						ArgObj := Variable {}
						ArgObj.TypeInfo = Type
						ArgObj.Name = Name
						ArgObj.Scope = Scope
						ArgObj.Internal = GenerateLocalIVN(&FObj, Type)
						ArgObj.Kind = VARIABLE

						FObj.Parameters = append(FObj.Parameters, ArgObj)
						DeclAppend(ArgObj)
					case shared.TokRParen:
						exit = true
					}
					if exit == true {
						break
					}
				}

				expect(shared.TokRParen)

				if peek(0).Value == "__attribute__" {
					expect(shared.TokIdent)
					expect(shared.TokLParen)
					expect(shared.TokLParen)

					exit := false
					for i < len(Tokens) {
						switch peek(0).Type {
						case shared.TokIdent:
							Attr := expect(shared.TokIdent)
							switch Attr {
							case "noreturn":
							default:
								error.Warning(11, "'" + Attr + "'", peek(-1), &Tokens)
							}
							FObj.Attributes = append(FObj.Attributes, Attr)
							if peek(0).Type != shared.TokRParen {
								expect(shared.TokComma)
							}
						case shared.TokRParen:
							exit = true
						}
						if exit == true {
							break
						}
					}

					expect(shared.TokRParen)
					expect(shared.TokRParen)
				}

				DeclAppend(FObj)
				Location := len((*TU).Declarations) - 1
		
				switch peek(0).Type {
				case shared.TokSemi:
					FObj.TypeInfo.Extern = true // implicit extern
					expect(shared.TokSemi)
					(*TU).Declarations[Location] = FObj
				default:
					expect(shared.TokLCurly)

					slice := []shared.Token {}

					depth := 1
					exit = false
					for j := i; j < len(Tokens); j++ {
						switch Tokens[j].Type {
						case shared.TokLCurly:
							depth++
							slice = append(slice, Tokens[j])
						case shared.TokRCurly:
							depth--
							if depth == 0 {
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

					CurrentFunction = &FObj
					ParseLocal(0, len(slice) - 1, Scope, slice, &FObj.Children, TU, false, false)

					(*TU).Declarations[Location] = FObj

					expect(shared.TokRCurly)
				}	
			case shared.TokEqual, shared.TokSemi:
				VObj := Variable {}
				VObj.Name = Name
				VObj.TypeInfo = TypeInformation
				VObj.Kind = VARIABLE
				VObj.Scope = Scope
				VObj.Internal = fmt.Sprintf("var_%d", IDCounter)
				IDCounter++

				if Scope != 0 {
					VObj.Internal =	GenerateLocalIVN(EnclosingFunction, TypeInformation)
				}

				DeclAppend(VObj)
				Location := len((*TU).Declarations) - 1

				if TypeInformation.Type == VOID && TypeInformation.PointerLength <= 0 {
					error.Error(7, "'void'", TypeTok, &Tokens)
				}

				if TypeInformation.Extern == true {	
					expect(shared.TokSemi)
					continue
				}

				switch peek(0).Type {
				case shared.TokSemi:
					expect(shared.TokSemi)
					continue
				default:
					expect(shared.TokEqual)
				}

				sloc := i
			
				exit := false
				end := 0
				for j := i; j < len(Tokens); j++ {
					switch Tokens[j].Type {
					case shared.TokSemi:
						end = j - i
						i = j
						exit = true
					}
					if exit == true {
						break
					}
				}

				end2 := i

				expect(shared.TokSemi)
	
				switch Scope {
				case 0:
					ParseLocal(sloc, end, Scope, Tokens, &VObj.Children, TU, true, true)
					(*TU).Declarations[Location] = VObj
				default:
					ParseLocal(NameLocation, end2, Scope, Tokens, VInitChildrenSlice, TU, false, false)	
					(*TU).Declarations[Location] = VObj
				}
			}
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

			DeclAppend(AsmObj)
		case shared.TokTypedef:
			expect(shared.TokTypedef)
			
			switch peek(0).Type {
			case shared.TokStruct:
				expect(shared.TokStruct)
				expect(shared.TokLCurly)

				MajorType := CompositeType {
					Type: STRUCT,
				}

				CurrentSize := 0

				for peek(0).Type != shared.TokRCurly {
					Type, _ := ParseType()
					CurrentSize += Type.Size
					Name := expect(shared.TokIdent)
					expect(shared.TokSemi)

					MemberObject := Type
					MemberObject.MemberName = Name

					MajorType.Children = append(MajorType.Children, MemberObject)
				}
				expect(shared.TokRCurly)

				MajorType.HighName = expect(shared.TokIdent)
				MajorType.Size = CurrentSize
				expect(shared.TokSemi)

				TypeMap = append(TypeMap, MajorType)
			default:
				Type, _ := ParseType()		
				Name := expect(shared.TokIdent)

				expect(shared.TokSemi)

				RealType := Type
				RealType.HighName = Name
				RealType.Size = Type.Size

				TypeMap = append(TypeMap, RealType)
			}
		}
	}	
}
