package neoparser

import (
	"lcc1/shared"
	"lcc1/error"
	"fmt"
)

var IDCounter = 1

func Parse(Tokens []shared.Token) *AST {
	TranslationUnit := AST {}

	ParseTop(Tokens, 0, &TranslationUnit)

	return &TranslationUnit
}

func ParseTop(Tokens []shared.Token, Scope int, TU *AST) { // TODO: add dynamic scoping for DECL	
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

	ParseType := func() CompositeType {
		Preset := 0
		PointerLength := 0

		Type := CompositeType {}

		exit := false
		for i < len(Tokens) {
			switch peek(0).Type {
			case shared.TokType:
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
						case 16:
							Type.Type = I16
						case 32:
							Type.Type = I32
						}
					}
				case "void":
					Type.Type = VOID
				case "char":
					Type.Type = I8
				}
			case shared.TokStar:
				expect(shared.TokStar)
				PointerLength++
			case shared.TokAmpersand:
				expect(shared.TokAmpersand)
				PointerLength--
			default:
				exit = true
			}
			if exit == true {
				break
			}
		}

		Type.PointerLength = PointerLength
		
		Size := 0

		if Type.PointerLength > 0 {
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

		return Type
	}

	for i < len(Tokens) {
		TypeInformation := ParseType()
		Name := expect(shared.TokIdent)

		switch peek(0).Type {
		case shared.TokLParen:
			FObj := Variable {}
			FObj.Name = Name
			FObj.Internal = Name
			FObj.TypeInfo = TypeInformation
			FObj.Kind = FUNCTION

			expect(shared.TokLParen)
			expect(shared.TokRParen)

			expect(shared.TokLCurly)

			slice := []shared.Token {}

			depth := 1
			exit := false
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


			Scope := CreateScope(0)
			ParseLocal(0, Scope, slice, &FObj, TU, false)

			expect(shared.TokRCurly)
			
			(*TU).Declarations = append((*TU).Declarations, FObj)
		case shared.TokIdent: // __attribute__
			fallthrough
		case shared.TokEqual:
			expect(shared.TokEqual)
		
			slice := []shared.Token {}

			exit := false
			for j := i; j < len(Tokens); j++ {
				switch Tokens[j].Type {
				case shared.TokSemi:
					i = j
					exit = true
				default:
					slice = append(slice, Tokens[j])
				}
				if exit == true {
					break
				}
			}

			expect(shared.TokSemi)
			
			VObj := Variable {}
			VObj.Name = Name
			VObj.TypeInfo = TypeInformation
			VObj.Internal = "var_" + fmt.Sprintf("%d", IDCounter)
			VObj.Kind = VARIABLE

			IDCounter++

			ParseLocal(0, Scope, slice, &VObj, TU, true)

			(*TU).Declarations = append((*TU).Declarations, VObj)
		}
	}	
}
