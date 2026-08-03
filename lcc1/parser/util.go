package parser

import (
	"lcc1/shared"
	"lcc1/error"
	"math"
	"strings"
	"strconv"
)

// TODO: migrate to strings.Builder!!!

func Write(text string, spaced bool) {
	if spaced == false {
		Code2 = Code2 + text + "\n"
	} else {
		Code2 = Code2 + "    " + text + "\n"
	}
}

func WritePre(text string, spaced bool) {
	if spaced == false {
		Code1 = Code1 + text + "\n"
	} else {
		Code1 = Code1 + "    " + text + "\n"
	}
}

func PreWrite(text string, spaced bool) {
	if spaced == false {
		Code1 = text + "\n" + Code1
	} else {
		Code1 = "    " + text + "\n" + Code1
	}
}

func CheckNum(token shared.Token) bool {
	if _, err := strconv.ParseInt(token.Value, 0, 64); err == nil {
		return true
	} else {
		return false
	}
}

func CreateStatic(variable Variable_Static) {
	WritePre(variable.Name + ":\n    .asciz \"" + variable.Value.(string) + "\"", false)	
}

func LookupParent(Scope int) int {
	if Scope == 1 {
		return -1
	}
	for _, s := range Scopes {
		if s.ID == Scope {
			return s.Parent
		}
	}
	return 1
}

func CreateScope(Parent int) int {
	Scopes = append(Scopes, Scope{ID: IDCounter, Parent: Parent})
	IDCounter++
	return IDCounter - 1
}

type GlobalInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func ReturnIntType(i int) string {
	switch {	
	case i <= math.MaxUint8 && i >= 0:
		return "unsigned short short int"
	case i <= math.MaxInt8 && i >= math.MinInt8:
		return "signed short short int"	
	case i <= math.MaxUint16 && i >= 0:
		return "unsigned short int"
	case i <= math.MaxInt16 && i >= math.MinInt16:
		return "signed short int"	
	case i <= math.MaxUint32 && i >= 0:
		return "unsigned long int"
	case i <= math.MaxInt32 && i >= math.MinInt32:
		return "signed long int"
	}
	return "unsigned short int"
}

func TypeToString(Type int, Pointer bool, PointerLength int) string {
	out := ""

	switch Type {
	case NUMBER8:
		out += "short short int"
	case NUMBER16:
		out += "int"
	case NUMBER32:
		out += "long int"
	case STRING:
		out += "char"
	case NULL:
		out += "void"
	case STRUCT:
		out += "struct"
	}

	for i := 0; i < PointerLength; i++ {
		out += "*"
	}

	return out
}

func ParseExpyL1(tokens []shared.Token, i int, Scope int) int {
	for {
		if i >= len(tokens) {
			break
		}
		i = ParseExpy(tokens, i, Scope, "r4", ArgumentTypeManifestEntry{Type: 999})
	}
	return i
}

func LookupType(name string) (TypeMapEntry, bool) {
	for _, Type := range TypeMap {
		if Type.Name == name {
			return Type, true
		}
	}
	return TypeMapEntry{ }, false
}

func ReturnStruct(name string) Variable_Static {
	for _, Variable := range Variables {
		if Variable.Name == name {
			return Variable
		}
	}
	return Variable_Static {}
}

func ParseNumberExpyDirect(tokens []shared.Token, i int, Scope int) (int, int) {
	expect := func(toktype shared.TokenType) string {
		var value string
		if i >= len(tokens) {
			if toktype != shared.TokSemi {
				error.Error(1, "'<EOF>'", tokens[i - 1], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
			return ""
		}
		if tokens[i].Type == toktype {
			value = tokens[i].Value
			i++
		} else {
			if toktype != shared.TokSemi && tokens[i].Type != shared.TokIdent {
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			} else if toktype == shared.TokSemi {
				error.Error(18, "", tokens[i - 1], &tokens)
			} else if tokens[i].Type == shared.TokIdent {
				error.Error(35, "", tokens[i], &tokens)
			}
			i++
		}
		return value
	}	
	peek := func(lookahead int) shared.Token {
		if i + lookahead < len(tokens) && i + lookahead >= 0 {
			return tokens[i + lookahead]
		}
		return shared.Token{Type: shared.TokEOF, Value: ""}
	}
	
	res := 0

	for {
		if i >= len(tokens) {
			break
		}
		
		exit := false
		exit_nodet := false

		switch peek(0).Type {
		case shared.TokAmpersand:
			expect(shared.TokAmpersand)
			label := expect(shared.TokIdent)
			Variable := LookupVariable(label, true, Scope, peek(-1), &tokens)
			WritePre(".ptr " + Variable.Real, true)
			return -1, i
		}
		num := expect(shared.TokNumber)

		OP_TRY:

		switch peek(0).Type {
		case shared.TokPlus, shared.TokMinus, shared.TokStar, shared.TokSlash:
		default:
			exit = true
		}
		if exit == true {
			if exit_nodet == false {
				n1_real, _ := strconv.ParseInt(num, 0, 64)
				res = int(n1_real)
			}
			break
		}
		
		op := peek(0).Value	
		expect(peek(0).Type)
		num2 := expect(shared.TokNumber)

		n1_real, _ := strconv.ParseInt(num, 0, 64)
		n2_real, _ := strconv.ParseInt(num2, 0, 64)

		switch op {
		case "+":
			res = int(n1_real) + int(n2_real)
		case "-":
			res = int(n1_real) - int(n2_real)
		case "*":
			res = int(n1_real) * int(n2_real)
		case "/":
			res = int(n1_real) / int(n2_real)
		}

		exit_nodet = true 
		goto OP_TRY
	}

	return res, i 
}

func LookupVariable(Name string, Enforce bool, Scope int, Token shared.Token, Tokens *[]shared.Token) Variable_Static {
	for {
		for _, variable := range Variables {
			if variable.Name == Name && variable.Scope == Scope {	
				return variable
			}
		}
		parent := LookupParent(Scope)
		if parent == -1 {	
			break
		}	
		Scope = parent
	}
	
	if Enforce == true {
		error.Error(4, "'" + Name + "'", Token, Tokens)
	}
	return Variable_Static{Name: "__ZERO", Type: NULL, Value: 0}
}

func LookupVariableDirect(Name string, Scope int) *Variable_Static {
	for {
		for i := 0; i < len(Variables); i++ {
			if Variables[i].Name == Name && Variables[i].Scope == Scope {	
				return &Variables[i]
			}
		}
		parent := LookupParent(Scope)
		if parent == -1 {	
			break
		}	
		Scope = parent
	}

	error.InternalCompilerError("Couldn't find variable on direct lookup")
	return &Variable_Static{Name: "__ZERO", Type: NULL, Value: 0}
}

func StringParse(tokens []shared.Token, start int) (string, int) {
	// Start would be the first token
	var str string = ""
	var loc int = 0
	
	if strings.HasPrefix(tokens[start].Value, "\"") == false {
		error.Error(2, "\"", tokens[start], &tokens)		
	}
	if strings.HasSuffix(tokens[start].Value, "\"") {
		word := strings.Trim(tokens[start].Value, "\"")
		str = word
		loc = start
	} else {
		var strtokens = []string { tokens[start].Value }
		for k := start + 1; k < len(tokens); k++ {
			strtokens = append(strtokens, tokens[k].Value)
			if strings.HasSuffix(tokens[k].Value, "\"") {
				start = k
				break
			}
		}
		str = strings.Join(strtokens, " ")
		str = strings.TrimSuffix(str,  "\"")
		loc = start
	}
	
	return str, loc
}

func FuncDeclLookup(Name string) (shared.Token, *[]shared.Token) {
	for _, d := range FunctionDecls {
		if d.Name == Name {
			return d.Token, &d.Set
		}
	}
	return shared.Token{Type: shared.TokEOF}, &[]shared.Token {}
}
