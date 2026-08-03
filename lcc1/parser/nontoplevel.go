package parser

import (
	"lcc1/error"
	"lcc1/shared"
	"strings"
	"fmt"
)

var CMP_OP string = ""
var _CMP_MOP_REVERSE string = ""
var _CMP_MOP string = ""
var _BREAK_TOPLEVEL string = ""
var _CONTINUE_TOPLEVEL string = ""
var _COERCE_TYPE int = 6
var _COERCE_PTR bool = false
var _COERCE_LENGTH int = 0
var _CURRENT_INFUNCTION Variable_Static

func ParseExpy(tokens []shared.Token, start int, Scope int, register string, RequiredType ArgumentTypeManifestEntry) int {
	i := start
	// CMP := false
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
			if toktype != shared.TokSemi {
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
		}
		return value
	}	
	peek := func(lookahead int) shared.Token {
		if i + lookahead < len(tokens) && i + lookahead >= 0 {
			return tokens[i + lookahead]
		}
		return shared.Token{Type: shared.TokEOF, Value: ""}
	}


	IDENT_FUNC := func(label string) {
		expect(shared.TokLParen)

		switch label {
		case "asm", "__asm__":
			if peek(0).Type == shared.TokQualifier && peek(0).Value == "volatile" {
				expect(shared.TokQualifier)
			}
			asmval := expect(shared.TokIdent)
			expect(shared.TokRParen)
			// expect(shared.TokSemi)
			// TODO: make quotes check here
			Write(strings.ReplaceAll(asmval, "\"", ""), true)
		case "sizeof":
			val := 0
			_label := expect(shared.TokIdent)
			expect(shared.TokRParen)
			Variable := LookupVariable(_label, true, Scope, peek(-2), &tokens)
			switch Variable.Type {
			case NUMBER8:
				val = 1
			case STRING:
				if Variable.Pointer == true {
					switch shared.Bits {
					case 16:
						val = 2
					case 32:
						val = 4
					}
				} else {
					val = 1
				}
			case NUMBER16:
				val = 2
			case NUMBER32:
				val = 4
			}
			Write("mov " + register + ", " + fmt.Sprintf("%d", val), true)
		default:
			NF_NOPARSE := false
			Function_Variable := LookupVariable(label, false, Scope, peek(-2), &tokens)
			if Function_Variable.Name == "__ZERO" {
				// Do a fallback to prevent a panic from Go
				NF_NOPARSE = true
				for i := 0; i < 6; i++ {
					Function_Variable.ArgumentTypeManifest = append(Function_Variable.ArgumentTypeManifest, ArgumentTypeManifestEntry{
						Type: 999,
					})
				}
			}
			// Parse arguments
			depth := 1
			pushed := 0
			j := i
			exit := false
			var CurrentTokens []shared.Token

			for j = i; j < len(tokens); j++ {
				if exit == true {
					break
				}

				switch tokens[j].Type {
				case shared.TokComma:
					if depth == 1 {
						ATMEntry := ArgumentTypeManifestEntry{ Type: 999 }	
						if pushed < Function_Variable.ArgNum {
							ATMEntry = Function_Variable.ArgumentTypeManifest[pushed]
						}

						ParseExpy(CurrentTokens, 0, Scope, "r7", ATMEntry)
						Write("push r7", true)
						CurrentTokens = []shared.Token{}
						pushed++
					} else {
						CurrentTokens = append(CurrentTokens, tokens[j])
					}
				case shared.TokLParen:
					depth++
					CurrentTokens = append(CurrentTokens, tokens[j])
				case shared.TokRParen:
					depth--
					if depth == 0 {
						exit = true
						break
					} else {
						CurrentTokens = append(CurrentTokens, tokens[j])
					}
				default:
					CurrentTokens = append(CurrentTokens, tokens[j])	
				}
			}

			// Push last args 
			if len(CurrentTokens) > 0 {
				ATMEntry := ArgumentTypeManifestEntry{ Type: 999 }	
				if pushed < Function_Variable.ArgNum {
					ATMEntry = Function_Variable.ArgumentTypeManifest[pushed]
				}

				ParseExpy(CurrentTokens, 0, Scope, "r7", ATMEntry)
				Write("push r7", true)
				pushed++
			}
			Write("call " + label, true)
			Write("mov " + register + ", e6", true)

			if NF_NOPARSE == false {
				if pushed < Function_Variable.ArgNum {
					t, s := FuncDeclLookup(label)	
					error.Error(20, "expected " + fmt.Sprintf("%d", Function_Variable.ArgNum) + ", have " + fmt.Sprintf("%d", pushed), peek(0), &tokens)
					error.Note(22, "'" + label + "' declared here", t, s)
				} else if pushed > Function_Variable.ArgNum {
					t, s := FuncDeclLookup(label)	
					error.Error(21, "expected " + fmt.Sprintf("%d", Function_Variable.ArgNum) + ", have " + fmt.Sprintf("%d", pushed), peek(0), &tokens)
					error.Note(22, "'" + label + "' declared here", t, s)
				}
			}
			i = j
		}
	}
	IDENT_STRING := func(label string) string {
		if label[len(label) - 1] != '"' {
			// TODO: fix "literal not terminated and have us handle it"
			error.Error(32, "\" character", peek(-1), &tokens)
		}
		_label := fmt.Sprintf("var_%d", IDCounter)
		IDCounter++
		__label := fmt.Sprintf("var_%d", IDCounter)
		IDCounter++

		WritePre(_label + ":", false)
		WritePre(".asciz \"" + strings.ReplaceAll(label, "\"", "") + "\"", true)
		WritePre(__label + ":", false)
		WritePre(".ptr " + _label, true)
		// Write("mov " + register + ", " + _label, true)

		return __label
	}
	_IDENT_INTENT := func(pointer bool, _type int, deref int, register bool, variable Variable_Static) string {	
		if peek(0).Type == shared.TokEqual || peek(0).Type == shared.TokIncrement || peek(0).Type == shared.TokDecrement {
			// Write intent (NEVER give one free dereference)
			Write("mov r2, r1", true)
			return "write"
		} else {
			// Read intent (ALWAYS give one free dereference)
			if deref >= 0 {
				if register == false {
					if pointer == false {
						switch _type {
						case NUMBER8, STRING, NULL:
							Write("lod r1, r2", true)
						case NUMBER16:
							Write("lod16 r1, r2", true)
						case NUMBER32:
							Write("lod32 r1, r2", true)
						case STRUCT:
							Write("mov r2, r1", true)
						}
					} else {
						Write("lod_ptr r1, r2", true)
					}
				} else {
					Write("mov r2, r1", true)
				}
			} else {
				if register == false {
					Write("mov r2, r1", true)
				} else {
					error.Error(37, "", peek(-1), &tokens)
				}
			}
		}
		return "read"
	}
	_NUMBER_PARSE := func(register string) {
		switch peek(0).Type {
		case shared.TokIdent, shared.TokStar, shared.TokAmpersand:
			subslice := []shared.Token {}
			fl_exit := false
			for {
				if fl_exit == true {
					break
				}
				switch peek(0).Type {
				case shared.TokStar, shared.TokAmpersand:
					subslice = append(subslice, peek(0))
					i++
				case shared.TokIdent:
					subslice = append(subslice, peek(0))
					i++
					fl_exit = true
				}
			}
			ParseExpy(subslice, 0, Scope, register, ArgumentTypeManifestEntry{Type: 999})
		case shared.TokNumber:

			Write("mov " + register + ", " + peek(0).Value, true)
			i++
		}
	}
	_PARSE_TYPE := func() (int, bool, int) {
		ptr := false
		long := false
		short := false
		shortshort := false
		unsigned := false
		constant := false
		bits := BitPref

		for {
			if i >= len(tokens) {
				break
			}
			if peek(0).Type == shared.TokQualifier {	
				qual := expect(shared.TokQualifier)	
				switch qual {
				case "short":
					if long == true {
						error.Error(12, "'long' declaration specifier", peek(-1), &tokens)
					}
					if short == true && shortshort == true {
						error.Warning(13, "'short' declaration specifier", peek(-1), &tokens)
					} else if short == true && shortshort == false {
						shortshort = true
						bits = 8
					} else {
						short = true
						bits = 16
					}
				case "long":
					if short == true {
						error.Error(12, "'short' declaration specifier", peek(-1), &tokens)
					}
					if long == true {
						error.Warning(13, "'long' declaration specifier", peek(-1), &tokens)
					}
					long = true
					bits = 32
				case "unsigned":
					if unsigned == true {
						error.Error(28, "'unsigned'", peek(-1), &tokens)
					}
					unsigned = true
				case "const":
					if constant == true {
						error.Error(28, "'const'", peek(-1), &tokens)
					}
					constant = true
				}
			} else {
				_type := expect(shared.TokType)
				var rtype int
				switch _type {
				case "int":
					switch bits {
					case 8:
						rtype = NUMBER8
					case 16:
						rtype = NUMBER16
					case 32:
						rtype = NUMBER32
					default:
						rtype = NUMBER16
					}
				case "char":
					if long == true || short == true || unsigned == true {
						error.Error(14, "for type 'char'", peek(-2), &tokens)
					}
					rtype = STRING
				case "void":
					rtype = NULL
				}

				ptr_length := 0
			_ptrtop:
				if peek(0).Type == shared.TokStar {
					ptr = true
					i++
					ptr_length++
					goto _ptrtop
				} else {
					if rtype == NULL && ptr_length < 1 {
						error.Error(7, "'void'", peek(-1), &tokens)	
					}
				}

				return rtype, ptr, ptr_length
			}
		}
		return NUMBER16, false, 0
	}
	_CMPOP_CLEANUP := func() {
		CMP_OP = ""
		_CMP_MOP_REVERSE = ""
		_CMP_MOP = ""
	}
	CheckRequiredType := func(CTYPE int, CPTR bool, CLENGTH int) {
		if RequiredType.Type != 999 {
			var RTYPE int
			var RPTR bool
			var RLENGTH int
			
			if RequiredType.Pointer == false {
				RTYPE = RequiredType.Type
			} else {
				RTYPE = RequiredType.Type2
			}
			RPTR = RequiredType.Pointer
			RLENGTH = RequiredType.PointerLength

			if CTYPE != RTYPE || CPTR != RPTR || CLENGTH != RLENGTH {
				error.Error(5, "passing '" + TypeToString(CTYPE, CPTR, CLENGTH) + "' to type of '" + TypeToString(RTYPE, RPTR, RLENGTH) + "'", peek(-1), &tokens)
			}	
		}
	}
	_TYPE_QUAL_DEFINE := func() {
		bodytokens := []shared.Token {}
		for j := i; j < len(tokens); j++ {
			bodytokens = append(bodytokens, tokens[j])
			if tokens[j].Type == shared.TokSemi {
				i = j + 1
				break
			}
		}
		level = 0
		L1_ALLOW_NONCONST = true
		IS_LOCAL = true
		Parse(bodytokens, Scope)
		L1_ALLOW_NONCONST = true
		IS_LOCAL = false
		level = 1
	}	
	_STRUCT_ACCESS := func(variable Variable_Static, CTYPE int, CPTR bool, CLENGTH int, MoveAfter bool, LoadArrow bool) (Variable_Static, int, bool, int) {
		FirstTime := true
		PointingToStruct := ""
		SA_TOP:
		if peek(0).Type == shared.TokPeriod || peek(0).Type == shared.TokArrow {
			switch peek(0).Type {
			case shared.TokPeriod:
				expect(shared.TokPeriod)
			case shared.TokArrow:
				expect(shared.TokArrow)
				if LoadArrow == true {
					Write("lod_ptr r1, r1", true)
				}
			}

			NM_NOWARN := false
			SNF_NOWARN := false

			if variable.Type != STRUCT && variable.Name != "__ZERO" && peek(-1).Type != shared.TokArrow {
				NM_NOWARN = true
				error.Error(50, "'" + TypeToString(variable.Type, variable.Pointer, variable.PointerLength) + "' is not a structure or union", peek(-1), &tokens)
			}
			if variable.Name == "__ZERO" {
				NM_NOWARN = true
			}
			if peek(-1).Type == shared.TokArrow && variable.Pointer != true {
				NM_NOWARN = true
				SNF_NOWARN = true
				error.Error(51, "'" + variable.Name + "' is not a pointer; did you mean to use '.'?", peek(-1), &tokens)
			}

			if peek(-1).Type == shared.TokArrow {
				if FirstTime == true {
					PointingToStruct = variable.PointingToStruct
				}
				Struct := ReturnStruct(PointingToStruct)
				if Struct.Name == "" && SNF_NOWARN == false {
					error.InternalCompilerError("struct not found on lookup using name '" + variable.PointingToStruct + "'")
				}

				variable = Struct
			}

			StructName := variable.Name
			identifier := expect(shared.TokIdent)
			found := false
			for _, Member := range variable.StructMemberList {
				if Member.Name == identifier {
					found = true
					Write("mov e8, " + fmt.Sprintf("%d", Member.Offset), true)
					Write("add r1, r1, e8", true)
					RequiredType = Member.RequiredType

					if Member.PointingToStruct != "" {
						PointingToStruct = Member.PointingToStruct
					}

					if Member.RequiredType.Pointer == false {
						CTYPE = Member.RequiredType.Type
					} else {
						CTYPE = Member.RequiredType.Type2
					}
					CPTR = Member.RequiredType.Pointer
					CLENGTH = Member.RequiredType.PointerLength

					variable.PointerLength = CLENGTH
					variable.Pointer = CPTR
					variable.Type = CTYPE
					break
				}
			}
			if found == false && NM_NOWARN == false {
				error.Error(49, "'" + identifier + "' in '" + StructName + "'", peek(-1), &tokens)
			}

			if MoveAfter == true {
				Write("mov " + register + ", r1", true)
			}
			
			if peek(0).Type == shared.TokArrow {
				FirstTime = false
				goto SA_TOP
			}

			return variable, CTYPE, CPTR, CLENGTH
		}
		return variable, CTYPE, CPTR, CLENGTH
	}

	deref := 0
	EQU_VT := NULL
	EQU_VAR := Variable_Static{}
	var op string = ""
	var NUM_TRY_DEREF bool
	var NUM_DIRECT bool

	EXPY_TOP:
	switch peek(0).Type {
	case shared.TokSemi:
		expect(shared.TokSemi)
		goto EXPY_TOP
	case shared.TokStar:
		expect(shared.TokStar)
		deref++
		goto EXPY_TOP
	case shared.TokAmpersand:
		expect(shared.TokAmpersand)
		deref--
		if deref <= -2 {
			error.Error(27, "'int'", peek(-1), &tokens)
		}
		goto EXPY_TOP
	case shared.TokGoto:
		// TODO: add goto var checks
		expect(shared.TokGoto)
		label := expect(shared.TokIdent)
		expect(shared.TokSemi)
		Write("jmp " + label, true)
		goto DONE
	case shared.TokType, shared.TokQualifier:
		_TYPE_QUAL_DEFINE()
		goto DONE
	case shared.TokLParen:
		expect(shared.TokLParen)
		_COERCE_TYPE, _COERCE_PTR, _COERCE_LENGTH = _PARSE_TYPE()	
		expect(shared.TokRParen)
		goto EXPY_TOP
	case shared.TokIf:
		// TODO: implement quick ifs

		IfScope := CreateScope(Scope)
		ElseScope := CreateScope(Scope)

		expect(shared.TokIf)
		expect(shared.TokLParen)

		exp_tokens := []shared.Token {}

		depth := 1
		exit := false
		for _ = i; i < len(tokens); i++ {	
			switch tokens[i].Type {
			case shared.TokLParen:
				depth++
				exp_tokens = append(exp_tokens, tokens[i])
			case shared.TokRParen:
				depth--
				if depth == 0 {
					exit = true
					break
				} else {
					exp_tokens = append(exp_tokens, tokens[i])
				}
			default:
				exp_tokens = append(exp_tokens, tokens[i])
			}
			if exit == true {
				// fmt.Println(tokens[i].Value)
				break
			}
		}
		ParseExpy(exp_tokens, 0, Scope, "r11", ArgumentTypeManifestEntry{Type: 999}) // r12 and r5 clobbered
											   // r11 result

		expect(shared.TokRParen)

		if_label := fmt.Sprintf("if_stmt_%d", IDCounter)
		IDCounter++
		else_label := fmt.Sprintf("else_stmt_%d", IDCounter)
		IDCounter++
		after_label := fmt.Sprintf("after_stmt_%d", IDCounter)
		IDCounter++	

		if_tokens := []shared.Token {}
		else_tokens := []shared.Token {}
		
		expect(shared.TokLCurly)	
		j := i
		depth = 1
		exit = false
		for j = i; j < len(tokens); j++ {
			if exit == true {
				break
			}
			switch tokens[j].Type {
			case shared.TokLCurly:
				depth++
				if_tokens = append(if_tokens, tokens[j])
			case shared.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					if_tokens = append(if_tokens, tokens[j])
				}
			default:
				if_tokens = append(if_tokens, tokens[j])
			}
		}
		i = j - 1
		expect(shared.TokRCurly)

		cmop := ""
		cmopr := ""
		if _CMP_MOP == "" {
			cmop = "jnz"
		} else {
			cmop = _CMP_MOP
		}

		if _CMP_MOP_REVERSE == "" {
			cmopr = "jz"
		} else {
			cmopr = _CMP_MOP_REVERSE
		}

		if peek(0).Type != shared.TokElse {
			// Write everything
			Write(cmop + " r11, " + if_label, true)
			Write(cmopr + " r11, " + after_label, true)
			Write(if_label + ":", false)
			ParseExpyL1(if_tokens, 0, IfScope)
			Write("jmp " + after_label, true)
			Write(after_label + ":", false)
			_CMPOP_CLEANUP()
			goto DONE
		}
		
		expect(shared.TokElse)

		expect(shared.TokLCurly)	
		j = i
		depth = 1
		exit = false
		for j = i; j < len(tokens); j++ {
			if exit == true {
				break
			}
			switch tokens[j].Type {
			case shared.TokLCurly:
				depth++
				else_tokens = append(else_tokens, tokens[j])
			case shared.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					else_tokens = append(else_tokens, tokens[j])
				}
			default:
				else_tokens = append(else_tokens, tokens[j])
			}
		}
		i = j - 1
		expect(shared.TokRCurly)

		// Write everything
		Write(cmop + " r11, " + if_label, true)
		Write(cmopr + " r11, " + else_label, true)	
		Write(if_label + ":", false)
		ParseExpyL1(if_tokens, 0, IfScope)
		Write("jmp " + after_label, true)
		Write(else_label + ":", false)
		ParseExpyL1(else_tokens, 0, ElseScope)
		Write(after_label + ":", false)
		_CMPOP_CLEANUP()
		goto DONE
	case shared.TokWhile:
		expect(shared.TokWhile)
		expect(shared.TokLParen)
		
		subslice := []shared.Token {}

		depth := 1
		exit := false

		top_label := "while_stmt_" + fmt.Sprintf("%d", IDCounter) + "_check"	
		middle_label := "while_stmt_" + fmt.Sprintf("%d", IDCounter) + "_body"
		bottom_label := "while_stmt_" + fmt.Sprintf("%d", IDCounter) + "_after"

		IDCounter++

		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case shared.TokLParen:
				depth++
				subslice = append(subslice, peek(0))
			case shared.TokRParen:
				depth--
				if depth < 1 {
					exit = true
				} else {
					subslice = append(subslice, peek(0))
				}
			default:
				subslice = append(subslice, peek(0))
			}
			if exit == true {
				break
			}
		}
		expect(shared.TokRParen)

		// Write check portion
		
		Write(top_label + ":", false)
		ParseExpy(subslice, 0, Scope, "r11", ArgumentTypeManifestEntry{Type: 999})

		cmop := ""
		if _CMP_MOP == "" {
			cmop = "jnz"
		} else {
			cmop = _CMP_MOP
		}

		Write(cmop + " r11, " + middle_label, true)
		Write("jmp " + bottom_label, true)

		expect(shared.TokLCurly)
		
		subslice2 := []shared.Token {}
		
		depth = 1
		exit = false
		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case shared.TokLCurly:
				depth++
				subslice2 = append(subslice2, peek(0))
			case shared.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice2 = append(subslice2, peek(0))
				}
			default:
				subslice2 = append(subslice2, peek(0))
			}
			if exit == true {
				break
			}
		}
		expect(shared.TokRCurly)

		WScope := CreateScope(Scope)

		otln_br := _BREAK_TOPLEVEL
		_BREAK_TOPLEVEL = bottom_label
		otln_co := _CONTINUE_TOPLEVEL
		_CONTINUE_TOPLEVEL = middle_label

		Write(middle_label + ":", false)
		ParseExpyL1(subslice2, 0, WScope)
		Write("jmp " + top_label, true)
		Write(bottom_label + ":", false)

		_BREAK_TOPLEVEL = otln_br
		_CONTINUE_TOPLEVEL = otln_co

		_CMPOP_CLEANUP()
		goto DONE
	case shared.TokDo:
		expect(shared.TokDo)
		expect(shared.TokLCurly)

		top_label := "do_stmt_" + fmt.Sprintf("%d", IDCounter) + "_top"
		middle_label := "do_stmt_" + fmt.Sprintf("%d", IDCounter) + "_check"
		bottom_label := "do_stmt_" + fmt.Sprintf("%d", IDCounter) + "_after"

		IDCounter++
		
		depth := 1
		exit := false
		subslice := []shared.Token {}

		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case shared.TokLCurly:
				depth++
				subslice = append(subslice, peek(0))
			case shared.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice = append(subslice, peek(0))
				}
			default:
				subslice = append(subslice, peek(0))
			}
			if exit == true {
				break
			}
		}
		expect(shared.TokRCurly)
		expect(shared.TokWhile)
		expect(shared.TokLParen)
		
		subslice2 := []shared.Token {}
		exit = false
		depth = 1
		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case shared.TokLParen:
				depth++
				subslice2 = append(subslice2, peek(0))
			case shared.TokRParen:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice2 = append(subslice2, peek(0))
				}
			default:
				subslice2 = append(subslice2, peek(0))
			}
			if exit == true {
				break
			}
		}

		expect(shared.TokRParen)
		expect(shared.TokSemi)

		otln := _BREAK_TOPLEVEL
		_BREAK_TOPLEVEL = bottom_label
		otln_co := _CONTINUE_TOPLEVEL	
		_CONTINUE_TOPLEVEL = middle_label

		Write(top_label + ":", false)
		DScope := CreateScope(Scope)
		ParseExpyL1(subslice, 0, DScope)
		Write(middle_label + ":", false)
		ParseExpy(subslice2, 0, DScope, "r11", ArgumentTypeManifestEntry{Type: 999})

		cmop := ""
		if _CMP_MOP == "" {
			cmop = "jnz"
		} else {
			cmop = _CMP_MOP
		}

		Write(cmop + " r11, " + top_label, true)
		Write("jmp " + bottom_label, true)

		Write(bottom_label + ":", false)

		_BREAK_TOPLEVEL = otln
		_CONTINUE_TOPLEVEL = otln_co
		_CMPOP_CLEANUP()
		goto DONE
	case shared.TokFor:
		expect(shared.TokFor)
		expect(shared.TokLParen)

		top_label := "for_stmt_" + fmt.Sprintf("%d", IDCounter) + "_check"
		bottom_label := "for_stmt_" + fmt.Sprintf("%d", IDCounter) + "_after"

		IDCounter++
		
		subslice := []shared.Token {}
		exit := false
		for _ = i; i < len(tokens); i++ {
			if exit == true {
				break
			}
			switch peek(0).Type {
			case shared.TokSemi:
				subslice = append(subslice, peek(0))
				exit = true
			default:
				subslice = append(subslice, peek(0))
			}
		}
		
		subslice2 := []shared.Token {}
		exit = false
		switch peek(0).Type {
		case shared.TokSemi:
			subslice2 = append(subslice2, shared.Token{Type: shared.TokNumber, Value: "1", Line: peek(0).Line, File: peek(0).File})
			expect(shared.TokSemi)
			goto COND_DONE	
		}	

		for _ = i; i < len(tokens); i++ {
			if exit == true {
				break
			}
			switch peek(0).Type {
			case shared.TokSemi:
				subslice2 = append(subslice2, peek(0))
				exit = true
			default:
				subslice2 = append(subslice2, peek(0))
			}
		}

		COND_DONE:

		FScope := CreateScope(Scope)
		ParseExpyL1(subslice, 0, FScope) // Initialize variable

		Write(top_label + ":", false)
		ParseExpy(subslice2, 0, FScope, "r11", ArgumentTypeManifestEntry{Type: 999})

		cmopr := ""
		if _CMP_MOP_REVERSE == "" {
			cmopr = "jz"
		} else {
			cmopr = _CMP_MOP_REVERSE
		}

		Write(cmopr + " r11, " + bottom_label, true)

		subslice3 := []shared.Token {}
		exit = false
		depth := 1

		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case shared.TokLParen:
				depth++
				subslice3 = append(subslice3, peek(0))
			case shared.TokRParen:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice3 = append(subslice3, peek(0))
				}
			default:
				subslice3 = append(subslice3, peek(0))
			}
			if exit == true {
				break
			}
		}

		expect(shared.TokRParen)
		expect(shared.TokLCurly)

		subslice4 := []shared.Token {}
		depth = 1
		exit = false
		
		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case shared.TokLCurly:
				depth++
				subslice4 = append(subslice4, peek(0))
			case shared.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice4 = append(subslice4, peek(0))
				}
			default:
				subslice4 = append(subslice4, peek(0))
			}
			if exit == true {
				break
			}
		}
		subslice3 = append(subslice3, shared.Token{Type: shared.TokSemi, Value: ";", Line: peek(-1).Line, File: peek(-1).File})

		otln := _BREAK_TOPLEVEL
		_BREAK_TOPLEVEL = bottom_label
		otln_co := _CONTINUE_TOPLEVEL	
		_CONTINUE_TOPLEVEL = top_label

		ParseExpyL1(subslice4, 0, FScope)
		ParseExpyL1(subslice3, 0, FScope)
		Write("jmp " + top_label, true)
		Write(bottom_label + ":", false)

		_BREAK_TOPLEVEL = otln
		_CONTINUE_TOPLEVEL = otln_co
		
		_CMPOP_CLEANUP()
		expect(shared.TokRCurly)
		goto DONE
	case shared.TokContinue:
		expect(shared.TokContinue)
		if _CONTINUE_TOPLEVEL == "" {
			error.Error(39, "", peek(-1), &tokens)
		} else {
			Write("jmp " + _CONTINUE_TOPLEVEL, true)
		}
		expect(shared.TokSemi)
	case shared.TokBreak:
		expect(shared.TokBreak)
		if _BREAK_TOPLEVEL == "" {
			error.Error(38, "", peek(-1), &tokens)
		} else {
			Write("jmp " + _BREAK_TOPLEVEL, true)
		}
		expect(shared.TokSemi)
	case shared.TokReturn:
		expect(shared.TokReturn)

		switch peek(0).Type {
		case shared.TokSemi:		
		default:
			if _CURRENT_INFUNCTION.Pointer == false && _CURRENT_INFUNCTION.Type == NULL {
				error.Error(44, "'" + _CURRENT_INFUNCTION.Name + "' should not return a value", peek(-1), &tokens)	
			}
			

			i = ParseExpy(tokens, i, Scope, "e6", ArgumentTypeManifestEntry{
				Type: _CURRENT_INFUNCTION.Type,
				Type2: _CURRENT_INFUNCTION.Type2,
				Pointer: _CURRENT_INFUNCTION.Pointer,
				PointerLength: _CURRENT_INFUNCTION.PointerLength,
			})
		}	
		expect(shared.TokSemi)

		Write("pop e11", true)
		Write("pop fp", true)
		Write("ret", true)
		goto DONE
	case shared.TokIdent:
		found := false
		for _, Type := range TypeMap {
			if Type.Name == peek(0).Value {
				found = true
				break
			}
		}
		if found == true {	
			_TYPE_QUAL_DEFINE()
			goto DONE
		}

		label := expect(shared.TokIdent)
		var variable Variable_Static

		// NUM_VAR_OVERRIDE = true

		var CTYPE int
		var CPTR bool
		var CLENGTH int

		if label[0] == '"' {
			LTL := IDENT_STRING(label)
			Write("mov r1, " + LTL, true)
			_IDENT_INTENT(true, NUMBER16, 0, false, variable)
			Write("mov " + register + ", r2", true)
			goto CONTINUE
		}

		array := false
		switch peek(0).Type {
		case shared.TokLParen:
			// TODO: allow comma arbitration
			// TODO: refactor this

			variable = LookupVariable(peek(-1).Value, false, Scope, peek(-1), &tokens)
			if variable.Name == "__ZERO" && peek(-1).Value != "asm" && peek(-1).Value != "sizeof" {
				error.Error(19, "'" + label + "'; ISO C99 and later do not support implicit function declarations", peek(-1), &tokens)
			}

			if _COERCE_TYPE != 6 {
				CTYPE = _COERCE_TYPE
				CPTR = _COERCE_PTR
				CLENGTH = _COERCE_LENGTH
				_COERCE_TYPE = 6
				_COERCE_PTR = false
				_COERCE_LENGTH = 0
			} else {
				if variable.Pointer == false {
					CTYPE = variable.Type
					CPTR = variable.Pointer
					CLENGTH = variable.PointerLength
				} else {
					CTYPE = variable.Type2
					CPTR = variable.Pointer
					CLENGTH = variable.PointerLength
				}
			}

			IDENT_FUNC(label)

			switch peek(0).Type {
			case shared.TokPeriod:
				error.UnimplementedMessage("Direct access of structs from functions is not currently supported")
			case shared.TokArrow:
				Write("mov r1, " + register, true)
				variable, CTYPE, CPTR, CLENGTH = _STRUCT_ACCESS(variable, CTYPE, CPTR, CLENGTH, true, false)
				EQU_VAR = variable
				EQU_VT = CTYPE
				switch peek(0).Type {
				case shared.TokEqual, shared.TokIncrement, shared.TokDecrement:
				default:
					if CPTR == true {
						Write("lod_ptr " + register + ", " + register, true)
					} else {
						switch CTYPE {
						case NUMBER8, STRING, NULL:
							Write("lod " + register + ", " + register, true)
						case NUMBER16:
							Write("lod16 " + register + ", " + register, true)
						case NUMBER32:
							Write("lod32 " + register + ", " + register, true)
						}
					}
				}
			}
			CheckRequiredType(CTYPE, CPTR, CLENGTH)
			goto CONTINUE
		case shared.TokColon:
			Write(label + ":", false)
			expect(shared.TokColon)
			goto DONE	
		}	

		variable = LookupVariable(label, true, Scope, peek(-1), &tokens)

		Write("mov r1, " + variable.Real, true)

		if _COERCE_TYPE != 6 {
			CTYPE = _COERCE_TYPE
			CPTR = _COERCE_PTR
			CLENGTH = _COERCE_LENGTH
			_COERCE_TYPE = 6
			_COERCE_PTR = false
			_COERCE_LENGTH = 0
		} else {
			if variable.Pointer == false {
				CTYPE = variable.Type
				CPTR = variable.Pointer
				CLENGTH = variable.PointerLength
			} else {
				// If encountering problems then come back and add variable.Type = _COERCE_TYPE
				CTYPE = variable.Type2
				CPTR = variable.Pointer
				CLENGTH = variable.PointerLength
			}
		}

		variable, CTYPE, CPTR, CLENGTH = _STRUCT_ACCESS(variable, CTYPE, CPTR, CLENGTH, false, true)

		switch peek(0).Type {
		case shared.TokLBracket:
			if variable.ArgNum < 1 {
				error.Error(40, "", peek(0), &tokens)
			}
			array = true
			expect(shared.TokLBracket)

			subslice := []shared.Token {}
			depth := 1
			exit := false
			for _ = i; i < len(tokens); i++ {
				switch peek(0).Type {
				case shared.TokLBracket:
					depth++
					subslice = append(subslice, peek(0))
				case shared.TokRBracket:
					depth--
					if depth == 0 {
						exit = true
					} else {
						subslice = append(subslice, peek(0))
					}
				default:
					subslice = append(subslice, peek(0))
				}
				if exit == true {
					break
				}
			}
			expect(shared.TokRBracket)

			ParseExpy(subslice, 0, Scope, "e8", ArgumentTypeManifestEntry{Type: 999})
		}

		if variable.Type != STRUCT {
			EQU_VT = variable.Type
		} else {
			EQU_VT = CTYPE
		}
		EQU_VAR = variable	
		if array == true {
			switch variable.Type {
			case NUMBER8:
				Write("add r1, r1, e8", true)
			case NUMBER16:
				Write("mov e7, 2", true)
				Write("mul e8, e8, e7", true)
				Write("add r1, r1, e8", true)
			case NUMBER32:
				Write("mov e7, 4", true)
				Write("mul e8, e8, e7", true)
				Write("add r1, r1, e8", true)
			}
		}
	
		Intent := _IDENT_INTENT(variable.Pointer, variable.Type, deref, variable.Register, variable)

		if Intent == "write" && deref < 0 {
			error.Error(45, "", peek(-1), &tokens)
		}

		x_deref := deref
		derefs := 0
		for x_deref > 0 {
			if RequiredType.Type != 999 {
				if CLENGTH > 0 {
					CLENGTH--
					if CLENGTH == 0 {
						CPTR = false
					}
				} else {
					CPTR = false
				}
			}
			x_deref--
			if variable.Pointer == true {
				if (derefs > variable.PointerLength - 1 && Intent == "write") || (derefs >= variable.PointerLength - 1 && Intent == "read") {
					if PIE == true {
						Write("add r2, r2, e14", true)
					}
					switch variable.Type2 {
					case NUMBER8, STRING, NULL:
						Write("lod r2, r2", true)
					case NUMBER16:
						Write("lod16 r2, r2", true)
					case NUMBER32:
						Write("lod32 r2, r2", true)
					}
				} else {
					Write("lod_ptr r2, r2", true)
				}	
			} else {
				switch variable.Type {
				case NUMBER8, STRING, NULL:
					Write("lod r2, r2", true)
				case NUMBER16:
					Write("lod16 r2, r2", true)
				case NUMBER32:
					Write("lod32 r2, r2", true)
				}
			}
			derefs++
		}

		if deref < 0 {
			CPTR = true
			CLENGTH = -deref
		}

		CheckRequiredType(CTYPE, CPTR, CLENGTH)
		Write("mov " + register + ", r2", true)	
	case shared.TokNumber:
		// Parse expressions
		// Load it up into r4

		_NUMBER_PARSE("r1")
		switch peek(0).Type {
		case shared.TokPlus, shared.TokMinus, shared.TokStar, shared.TokSlash, shared.TokPercent:
			NUM_DIRECT = true
		default:
			NUM_DIRECT = true	
		}	
		NUM_TRY_DEREF = true	
	}

	switch peek(0).Type {
	default:
		if NUM_DIRECT == true {
			Write("mov " + register + ", r1", true)
		}
	case shared.TokPlus, shared.TokMinus, shared.TokStar, shared.TokSlash, shared.TokPercent:
		if NUM_TRY_DEREF == false {
			Write("mov r5, " + register, true)
		}

		if NUM_DIRECT == true {
			Write("mov r5, r1", true)
			NUM_DIRECT = false
		}

		OP_TRY:
		op = expect(peek(0).Type)
		_NUMBER_PARSE("r6")

		switch op {
		case "+":
			Write("add r1, r5, r6", true)
		case "-":
			Write("sub r1, r5, r6", true)
		case "*":	
			Write("mul r1, r5, r6", true)
		case "/":
			Write("div r1, r5, r6", true)
		case "%":
			Write("mod r1, r5, r6", true)
		case "&":
			Write("and r1, r5, r6", true)
		case "|":
			Write("or r1, r5, r6", true)
		case "^":
			Write("xor r1, r5, r6", true)
		case "<<":
			Write("shl r1, r5, r6", true)
		case ">>":
			Write("shr r1, r5, r6", true)
		}
		switch peek(0).Type {
		case shared.TokPlus, shared.TokMinus, shared.TokStar, shared.TokSlash:
			Write("mov r5, r1", true)
			goto OP_TRY
		}
		if NUM_TRY_DEREF == false {
			Write("mov " + register + ", r1", true)
		}
		NUM_DIRECT = false
	}

	if NUM_TRY_DEREF == true && _COERCE_PTR == true {
		var CPTR bool
		var CLENGTH int

		if _COERCE_TYPE != 6 {
			CPTR = _COERCE_PTR
			CLENGTH = _COERCE_LENGTH	
		}

		FakeVar := Variable_Static {
			Name: "__FAKE__",
			Pointer: true,
			PointerLength: _COERCE_LENGTH,
			Type2: _COERCE_TYPE,
		}

		Intent := _IDENT_INTENT(_COERCE_PTR, _COERCE_TYPE, deref, true, FakeVar)

		x_deref := deref
		derefs := 0
		for x_deref - 1 > 0 {
			if RequiredType.Type != 999 {
				if CLENGTH > 0 {
					CLENGTH--
					if CLENGTH == 0 {
						CPTR = false
					}
				} else {
					CPTR = false
				}
			}
			x_deref--
			if (derefs > _COERCE_LENGTH - 1 && Intent == "write") || (derefs >= _COERCE_LENGTH - 1 && Intent == "read") {
				if PIE == true {
					Write("add r2, r2, e14", true)
				}
				switch _COERCE_TYPE {
				case NUMBER8, STRING, NULL:
					Write("lod r2, r2", true)
				case NUMBER16:
					Write("lod16 r2, r2", true)
				case NUMBER32:
					Write("lod32 r2, r2", true)
				}
			} else {
				Write("lod_ptr r2, r2", true)
			}	
			derefs++
		}
	
		Write("mov " + register + ", r2", true)

		// Construct fake variable
		// Removed the pointer variant because causing issues.
		// ^ I added it back

		if RequiredType.Type != 999 {
			required := ""
			actual := ""

			if RequiredType.Pointer == false {
				required = "integer"
			} else {
				required = "pointer"
			}
			
			if CPTR == false {
				actual = "integer"
			} else {
				actual = "pointer"
			}

			if (RequiredType.Pointer != CPTR) || (RequiredType.PointerLength != CLENGTH) {
				error.Error(5, "passing '" + actual + "' to '" + required + "'", peek(0), &tokens)
			}
		}

		if _COERCE_PTR == false {
			EQU_VAR = Variable_Static {Pointer: _COERCE_PTR, Type: _COERCE_TYPE}
		} else {
			EQU_VAR = Variable_Static {Pointer: _COERCE_PTR, Type2: _COERCE_TYPE}
		}
		
		EQU_VT = _COERCE_TYPE
		_COERCE_TYPE = 6
		_COERCE_PTR = false
		_COERCE_LENGTH = 0
	}

	if NUM_TRY_DEREF == true && deref < 1 && peek(0).Type == shared.TokEqual {	
		error.Error(45, "", peek(0), &tokens)
	}

	_COERCE_TYPE = 6
	_COERCE_PTR = false
	_COERCE_LENGTH = 0
	NUM_TRY_DEREF = false

	
	CONTINUE:
	if peek(0).Type != shared.TokEqual && peek(0).Type != shared.TokEquality && peek(0).Type != shared.TokInequality && peek(0).Type != shared.TokGEqual && peek(0).Type != shared.TokLEqual && peek(0).Type != shared.TokLAngle && peek(0).Type != shared.TokRAngle && peek(0).Type != shared.TokIncrement && peek(0).Type != shared.TokDecrement {
		// expect(shared.TokSemi)
		return i
	}
	
	switch peek(0).Type {
	case shared.TokEquality, shared.TokInequality, shared.TokGEqual, shared.TokLEqual, shared.TokLAngle, shared.TokRAngle:
		cmpopreal := ""
		expect(peek(0).Type)
		switch peek(-1).Type {
		case shared.TokEquality:
			_CMP_MOP = "jnz"
			_CMP_MOP_REVERSE = "jz"
			if CMP_OP == "" {
				cmpopreal = "cmp"
			} else {
				cmpopreal = CMP_OP
			}
		case shared.TokInequality:
			_CMP_MOP = "jz"
			_CMP_MOP_REVERSE = "jnz"
			if CMP_OP == "" {
				cmpopreal = "cmp"
			} else {
				cmpopreal = CMP_OP
			}
		case shared.TokGEqual:
			cmpopreal = "ilt"
			_CMP_MOP = "jz"
			_CMP_MOP_REVERSE = "jnz"
		case shared.TokLEqual:
			cmpopreal = "igt"
			_CMP_MOP = "jz"
			_CMP_MOP_REVERSE = "jnz"
		case shared.TokLAngle:
			cmpopreal = "ilt"
			_CMP_MOP = "jnz"
			_CMP_MOP_REVERSE = "jz"
		case shared.TokRAngle:
			cmpopreal = "igt"
			_CMP_MOP = "jnz"
			_CMP_MOP_REVERSE = "jz"
		}

		i = ParseExpy(tokens, i, Scope, "r5", ArgumentTypeManifestEntry{Type: 999})
		// Universal IF register = r11
	
		Write(cmpopreal + " r11, " + register + ", r5", true)	
	default:
		switch peek(0).Type {
		default:
			expect(shared.TokEqual)
			i = ParseExpy(tokens, i, Scope, "r5", ArgumentTypeManifestEntry{Type: 999})
		case shared.TokIncrement, shared.TokDecrement:
			ParseExpy([]shared.Token{
				shared.Token {
					Type: shared.TokIdent,
					Value: EQU_VAR.Name,	
				},
			}, 0, Scope, "r0", ArgumentTypeManifestEntry{Type: 999})

			Write("mov r3, " + register, true)

			Write("mov " + register + ", r0", true)

			Write("mov r5, r0", true);

			if peek(0).Type == shared.TokIncrement {
				Write("inc r0", true)
			} else {
				Write("dec r0", true)
			}

			if EQU_VAR.Pointer == false {	
				switch EQU_VT {
				case NUMBER8, STRING, NULL:
					Write("str r3, r0", true)
				case NUMBER16:
					Write("str16 r3, r0", true)	
				case NUMBER32:
					Write("str32 r3, r0", true)
				}
			} else {
				if deref >= EQU_VAR.PointerLength {
					switch EQU_VAR.Type2 {
					case NUMBER8, STRING, NULL:
						Write("str r3, r0", true)
					case NUMBER16:
						Write("str16 r3, r0", true)
					case NUMBER32:
						Write("str32 r3, r0", true)
					}
				} else {
					Write("str_ptr r3, r0", true)
				}
			}

			expect(peek(0).Type)
			return i
		}
		
		if EQU_VAR.Const == true {
			error.Error(33, "'" + EQU_VAR.Name + "' with const-qualified type", peek(-1), &tokens)
			token, stream := FuncDeclLookup(EQU_VAR.Name)
			error.Note(22, "'" + EQU_VAR.Name + "' declared here", token, stream)
		}

		if EQU_VAR.Pointer == false {
			switch EQU_VAR.Type {
			case STRUCT:
				Write("push " + register, true)
				Write("push r5", true)
				Write("push " + fmt.Sprintf("%d", EQU_VAR.StructTotalSize), true)
				funcname := "_builtin_lcc_memcpy"
				switch shared.Bits {
				case 32:
					funcname += "32"
				default:
					funcname += "16"
				}
				Write("call " + funcname, true)
			default:
				switch EQU_VT {
				case NUMBER8, STRING, NULL:
					Write("str " + register + ", r5", true)
				case NUMBER16:
					Write("str16 " + register + ", r5", true)	
				case NUMBER32:
					Write("str32 " + register + ", r5", true)
				}	
			}	
		} else {
			if deref >= EQU_VAR.PointerLength {
				switch EQU_VAR.Type2 {
				case NUMBER8, STRING, NULL:
					Write("str " + register + ", r5", true)
				case NUMBER16:
					Write("str16 " + register + ", r5", true)
				case NUMBER32:
					Write("str32 " + register + ", r5", true)
				}
			} else {
				Write("str_ptr " + register + ", r5", true)
			}
		}
	
		if peek(-1).Type == shared.TokEqual {
			expect(shared.TokSemi)
		}
	}

	DONE:
	_COERCE_TYPE = 6
	_COERCE_PTR = false
	_COERCE_LENGTH = 0
	return i
}
