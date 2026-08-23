package lexer

import (
	"strconv"
	"fmt"
	"lcc1/shared"
	"strings"
)

func Lex(code []SmallToken, filename string) []shared.Token {
	var tokens = []shared.Token {}

	Add := func(Type shared.TokenType, Value string, ST SmallToken) {
		MathToken := false

		switch Type {
		case shared.TokPlus, shared.TokMinus, shared.TokStar, shared.TokSlash, shared.TokPercent, shared.TokShiftLeft, shared.TokShiftRight, shared.TokIncrement, shared.TokDecrement:
			MathToken = true
		}

		tokens = append(tokens, shared.Token{
			Type: Type,
			Value: Value,
			FakeValue: Value,
			Line: ST.Line,
			File: ST.Filename,
			MathToken: MathToken,
		})
	}

	AddWithFake := func(Type shared.TokenType, Value string, RealValue string, ST SmallToken) {
		MathToken := false

		switch Type {
		case shared.TokPlus, shared.TokMinus, shared.TokStar, shared.TokSlash, shared.TokPercent, shared.TokShiftLeft, shared.TokShiftRight, shared.TokIncrement, shared.TokDecrement:
			MathToken = true
		}

		tokens = append(tokens, shared.Token{
			Type: Type,
			Value: RealValue,
			FakeValue: Value,
			Line: ST.Line,
			File: ST.Filename,
			MathToken: MathToken,
		})
	}
   	 
	for i := 0; i < len(code); i++ {
		SToken := code[i]
		content := SToken.Value

		switch content {	
		case "int", "void", "char":
			Add(shared.TokType, content, SToken)
		case "volatile", "unsigned", "short", "long", "static", "const", "extern":
			Add(shared.TokQualifier, content, SToken)
		case "return":
			Add(shared.TokReturn, content, SToken)
		case "if":
			Add(shared.TokIf, content, SToken)
		case "else":
			Add(shared.TokElse, content, SToken)
		case "break":
			Add(shared.TokBreak, content, SToken)
		case "continue":
			Add(shared.TokContinue, content, SToken)
		case "(":
			Add(shared.TokLParen, content, SToken)
		case ")":
			Add(shared.TokRParen, content, SToken)
		case "{":
			Add(shared.TokLCurly, content, SToken)
		case "}":
			Add(shared.TokRCurly, content, SToken)
		case ";":
			Add(shared.TokSemi, content, SToken)
		case "+":
			Add(shared.TokPlus, content, SToken)
		case "-":
			Add(shared.TokMinus, content, SToken)
		case "*":
			Add(shared.TokStar, content, SToken)
		case "/":	
			Add(shared.TokSlash, content, SToken)
		case "%":
			Add(shared.TokPercent, content, SToken)
		case "=":
			Add(shared.TokEqual, content, SToken)
		case "==":
			Add(shared.TokEquality, content, SToken)
		case ",":
			Add(shared.TokComma, content, SToken)
		case ":":
			Add(shared.TokColon, content, SToken)
		case "goto":
			Add(shared.TokGoto, content, SToken)
		case "for":
			Add(shared.TokFor, content, SToken)
		case "while":
			Add(shared.TokWhile, content, SToken)
		case "do":
			Add(shared.TokDo, content, SToken)
		case "<":
			Add(shared.TokLAngle, content, SToken)
		case "<=":
			Add(shared.TokLEqual, content, SToken)
		case ">":
			Add(shared.TokRAngle, content, SToken)
		case ">=":
			Add(shared.TokGEqual, content, SToken)
		case "&":
			Add(shared.TokAmpersand, content, SToken)
		case "!":
			Add(shared.TokExclamation, content, SToken)
		case "!=":
			Add(shared.TokInequality, content, SToken)
		case "//":
		case "\n":
		case "[":
			Add(shared.TokLBracket, content, SToken)
		case "]":
			Add(shared.TokRBracket, content, SToken)
		case "typedef":
			Add(shared.TokTypedef, content, SToken)
		case "struct":
			Add(shared.TokStruct, content, SToken)
		case "++":
			Add(shared.TokIncrement, content, SToken)
		case "--":
			Add(shared.TokDecrement, content, SToken)
		case ".":
			Add(shared.TokPeriod, content, SToken)
		case "->":
			Add(shared.TokArrow, content, SToken)
		case "<<":
			Add(shared.TokShiftLeft, content, SToken)
		case ">>":
			Add(shared.TokShiftRight, content, SToken)
		default:
			num, err := strconv.ParseInt(content, 0, 64)
			if err == nil {
				AddWithFake(shared.TokNumber, content, fmt.Sprintf("%d", num), SToken)
			} else {
				switch content[0] {
				case '\'':
					str := "0x"
					content2 := strings.ReplaceAll(content, "'", "")
					for _, b := range content2 {
						str += fmt.Sprintf("%02x", b)
					}
					AddWithFake(shared.TokNumber, content, str, SToken)
				default:
					Add(shared.TokIdent, content, SToken)
				}	
			}
		}	
	}

	return tokens
}
