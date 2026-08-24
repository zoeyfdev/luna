package parser

import (	
	"lcc1/shared"
	"lcc1/error"
	"fmt"
	"strings"
	"math"
	"strconv"
)

var topLevelName string
var BitPref int = 16
var IS_LOCAL bool = false
var CurrentDisplacement uint32

func Parse(tokens []shared.Token, Scope int) {
	error.WarningNoGaze(55, "", shared.Token{})
	i := 0
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
		} else {
			if toktype != shared.TokSemi {
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
		}
		i++
		return value
	}	
	peek := func(lookahead int) shared.Token {
		if i + lookahead < len(tokens) {
			return tokens[i + lookahead]
		}
		return shared.Token{Type: shared.TokEOF, Value: ""}
	}

	_PARSE_ATTR := func(name string) []string {
		var attrs []string
		var _RETURNS []string

		expect(shared.TokIdent)
		expect(shared.TokLParen)
		expect(shared.TokLParen)

		expComma := false
		for {
			if expComma == false {
				attr := expect(shared.TokIdent)
				attrs = append(attrs, attr)
				expComma = true
				if peek(0).Type == shared.TokRParen {
					break
				}
			} else {
				expect(shared.TokComma)
				expComma = false
			}
		}	
		
		expect(shared.TokRParen)
		expect(shared.TokRParen)

		for _, attr := range attrs {
			switch attr {
			case "norename":
				_RETURNS = append(_RETURNS, "norename")	
			case "noreturn":
				_RETURNS = append(_RETURNS, "noreturn")
			case "require_const":
				_RETURNS = append(_RETURNS, "require_const")
			default:
				error.Warning(11, "'" + attr + "'", tokens[i - 3], &tokens)
			}
		}
		return _RETURNS
	}

	_PARSE_TYPE := func() (int, bool, int) {
		ptrlen := 0
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
				_type := "int"
				switch peek(0).Type {
				case shared.TokType, shared.TokIdent:
					_type = expect(peek(0).Type)
				default:
					_type = expect(shared.TokType)
				}
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
				default:
					var TypeEntry TypeMapEntry
					found := false
					for _, Type := range TypeMap {
						if Type.Name == _type {
							found = true
							TypeEntry = Type
							break
						}
					}
					if found == false {
						error.Error(52, "'" + _type + "'", peek(-1), &tokens)	
					}
					rtype = TypeEntry.ReferralType
					if TypeEntry.Pointer == true {
						ptr = true
						ptrlen += TypeEntry.PointerLength
					}
				}

			ptrtop:
				if peek(0).Type == shared.TokStar {
					ptr = true
					i++
					ptrlen++
					goto ptrtop
				} else {
					if rtype == NULL && ptrlen < 1 {
						error.Error(7, "'void'", peek(-1), &tokens)	
					}
				}

				return rtype, ptr, ptrlen
			}
		}
		return NUMBER16, false, 0
	}

	ReturnDisplacement := func(Type int) uint32 {
		OldDisplacement := CurrentDisplacement
		switch Type {
		case NUMBER8:
			CurrentDisplacement += 1
		case NUMBER16:
			CurrentDisplacement += 2
		case NUMBER32:
			CurrentDisplacement += 4
		case STRING:
			CurrentDisplacement += 4 // Just to be safe
		case NULL:
			CurrentDisplacement += 4 // Just to be safe as well
		default:
			CurrentDisplacement += 4 // Just to be safe as well :3
		}
		return OldDisplacement
	}
	
	for {
		if i >= len(tokens) {
			break
		}	

		switch level {
		case 0:	
			var ptr bool = false
			_typetoken := tokens[i]

			long := false
			short := false
			shortshort := false
			unsigned := false
			constant := false
			extern := false
			static := false
			breakoff := false
			bits := BitPref	
			for {
				if i >= len(tokens) {
					breakoff = true
					break
				}
				if peek(0).Value == "asm" || peek(0).Value == "__asm__" {
					expect(shared.TokIdent)
					if peek(0).Value == "volatile" {
						expect(shared.TokQualifier)
					}
					expect(shared.TokLParen)
					str, end := StringParse(tokens, i)
					i = end + 1
					Write(str, false)
					expect(shared.TokRParen)
					expect(shared.TokSemi)
					continue
				}
				if peek(0).Value == "__embed__" {
					pre := false
					static := false
					name := ""

					expect(shared.TokIdent)

					arg_parse_top:
					switch peek(0).Value {
					case "pre":
						pre = true
						i++
						goto arg_parse_top
					case "static":
						static = true
						i++
						goto arg_parse_top
					}	

					name = expect(shared.TokIdent)
					// __embed__ pre SOUND_LABEL

					expect(shared.TokLParen)
					expect(shared.TokLParen)
					path := expect(shared.TokIdent)
					expect(shared.TokRParen)
					expect(shared.TokRParen)
					expect(shared.TokSemi)

					if pre == false {
						Write(name + ":", false)
						Write(".embed \"" + strings.ReplaceAll(path, "\"", "") + "\"", true)
					} else {
						WritePre(name + ":", false)
						WritePre(".embed \"" + strings.ReplaceAll(path, "\"", "") + "\"", true)
					}
					if static == false {
						PreWrite(".global " + name, false)
					}
					continue
				}

				if peek(0).Type == shared.TokTypedef {
					expect(shared.TokTypedef)
					
					switch peek(0).Type {
					case shared.TokQualifier, shared.TokType:
						Type, Pointer, PointerLength := _PARSE_TYPE()
						if Pointer == true {
							error.Warning(41, "", peek(-1), &tokens);
						}
						

						name := expect(shared.TokIdent)
						expect(shared.TokSemi)

						TypeMap = append(TypeMap, TypeMapEntry {
							Name: name,
							ReferralType: Type,
							Pointer: Pointer,
							PointerLength: PointerLength,

						})	
					case shared.TokStruct:
						expect(shared.TokStruct)
						expect(shared.TokLCurly)

						var MemberList []StructMemberListEntry	
						size_accumulator := 0;

						CheckMemberOfStruct := func(name string) bool {
							for _, Member := range MemberList {
								if Member.Name == name {
									return true
								}
							}
							return false
						}

						for {
							tname := peek(0).Value
							_, IsType := LookupType(peek(0).Value)
							if peek(0).Type == shared.TokQualifier || peek(0).Type == shared.TokType || IsType == true {
								Type, Pointer, PtrLen := _PARSE_TYPE()
								_size := 0

								if Type == STRUCT && Pointer == false {
									error.UnimplementedMessage("Directly including structs in structs is not currently supported")
								}
								
								if Pointer == false {
									switch Type {
									case NUMBER8, STRING, NULL:
										_size += 1
									case NUMBER16:
										_size += 2
									case NUMBER32:
										_size += 4
									}
								} else {
									switch shared.Bits {
									case 32:
										_size += 4
									default:
										_size += 2
									}
								}
								
								member_name := expect(shared.TokIdent)

								if CheckMemberOfStruct(member_name) == true {
									error.Error(48, "'" + member_name + "'", peek(-1), &tokens)
								}

								MemberList = append(MemberList, StructMemberListEntry {
									Name: member_name,
									Offset: size_accumulator,
									PointingToStruct: tname,
									RequiredType: ArgumentTypeManifestEntry {
										Type: Type,
										Pointer: Pointer,
										PointerLength: PtrLen,
									},
								})
								
								size_accumulator += _size
								expect(shared.TokSemi)
							} else {
								if peek(0).Type == shared.TokRCurly {
									i++
									break
								}
								error.Error(47, "", peek(0), &tokens)
								i++
							}
						}
						struct_name := expect(shared.TokIdent)
						expect(shared.TokSemi)

						StructVariable := Variable_Static {
							Name: struct_name,
							Type: STRUCT,
							StructTotalSize: size_accumulator,
							StructMemberList: MemberList,
						}

						Variables = append(Variables, StructVariable)	
						IDCounter++

						TypeMap = append(TypeMap, TypeMapEntry {
							ReferralType: STRUCT,
							EmbeddedStruct: StructVariable,
							Name: struct_name,
						})
					}
					continue
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
					case "extern":
						if extern == true {
							error.Error(28, "'extern'", peek(-1), &tokens)
						}
						extern = true
					case "static":
						if static == true {
							error.Error(28, "'static'", peek(-1), &tokens)
						}
						static = true
					}
				} else {
					break
				}
			}
			if breakoff == true {
				break
			}
	
			_type := ""
			_type_tok := peek(0)
			switch peek(0).Type {
			case shared.TokType, shared.TokIdent:
				_type = expect(peek(0).Type)
			default:
				_type = expect(shared.TokType)
			}
			_ptrlen := 0
		_ptrtop:
			if peek(0).Type == shared.TokStar {
				ptr = true
				i++
				_ptrlen++
				goto _ptrtop
			}

			// TODO: add ampersand as a reverser so we can do &* as a no-op

			name := expect(shared.TokIdent)
		
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
				if long == true || short == true || unsigned == true {
					error.Error(14, "for type 'void'", peek(-2), &tokens)
				}
				rtype = NULL
			default:
				found := false
				for _, Type := range TypeMap {
					if Type.Name == _type {
						rtype = Type.ReferralType
						found = true
					}
				}
				if found == false {
					error.Error(25, "", peek(-1), &tokens)
				}
			}	

			_variable := LookupVariable(name, false, Scope, tokens[i - 1], &tokens) 
			if _variable.Name != "__ZERO" && _variable.Scope == Scope {	
				error.Error(3, "'" + name + "'", tokens[i - 1], &tokens)
			}

			FunctionDecls = append(FunctionDecls, FunctionDecl {
				Name: name, 
				Token: peek(-1), 
				Set: tokens,
			})	

			var attrs []string
			var UnpackOrders []UnpackOrder
			var ManifestEntries []ArgumentTypeManifestEntry

			allow_nonconst := L1_ALLOW_NONCONST

			switch peek(0).Type {
			case shared.TokLParen:
				CurrentDisplacement = 0	

				if rtype == STRUCT && ptr != true {
					error.UnimplementedMessage("direct returning of structs is not supported due to ABI limitations")
				}

				expect(shared.TokLParen)
				fscope := CreateScope(Scope)	
		
				register := 0
				nargs := 0
				var BasinSize uint32
				switch peek(0).Type {
				case shared.TokType, shared.TokQualifier:	
					register = 0
					expComma := false
					for j := i; j < len(tokens); j++ {
						if peek(0).Type == shared.TokRParen {	
							expect(shared.TokRParen)
							break
						}
						if expComma == false {
							if register >= 6 {
								error.Error(9, "", peek(0), &tokens)	
							}

							__arg_reg := fmt.Sprintf("e%d", register)

							// __rn := fmt.Sprintf("var_%d", IDCounter)
							__rn := "fp + "

							IDCounter++
							__rtype, __ptr, __ptrlen := _PARSE_TYPE()
							__name := expect(shared.TokIdent)	
							
							if __ptr == false {
								ManifestEntries = append(ManifestEntries, ArgumentTypeManifestEntry{
									Type: __rtype,
									Pointer: false,
									PointerLength: __ptrlen,
								})
							} else {
								ManifestEntries = append(ManifestEntries, ArgumentTypeManifestEntry{
									Type2: __rtype,
									Pointer: true,
									PointerLength: __ptrlen,
								})
							}

							var dis uint32

							if extern == true {
								goto ARG_DECL_DONE
							}

							if __ptr == false {
								dis = ReturnDisplacement(__rtype)
							} else {
								dis = ReturnDisplacement(NUMBER32) 
							}
							
							__rn += fmt.Sprintf("%d", dis)
							
							UnpackOrders = append(UnpackOrders, UnpackOrder{Register: __arg_reg, Label: __rn, Type: __rtype, Pointer: __ptr, PointerLength: __ptrlen})

							switch __ptr {
							case false:
								switch __rtype {
								case NUMBER8:	
									BasinSize += 1
								case NUMBER16:
									BasinSize += 2
								case NUMBER32, STRING, NULL:
									BasinSize += 4
								}
							case true:
								BasinSize += 4
							}
							
							ARG_DECL_DONE:
							if __ptr == false {
								Variables = append(Variables, Variable_Static{Name: __name, Type: __rtype, Value: nil, Scope: fscope, Real: __rn, Pointer: __ptr})
							} else {
								Variables = append(Variables, Variable_Static{Name: __name, Type: NUMBER16, Type2: __rtype, Value: nil, Scope: fscope, Real: __rn, Pointer: __ptr, PointerLength: __ptrlen})
							}
							
							register++
							nargs++
							expComma = true
						} else {
							expect(shared.TokComma)
							expComma = false
						}	
					}	
				case shared.TokRParen:
					expect(shared.TokRParen)
				}

				PTS := ""
				if rtype == STRUCT && ptr == true {
					PTS = _type
				}

				var FuncVar Variable_Static
				if ptr == false {
					FuncVar = Variable_Static{
						Name: name, 
						Type: rtype, 
						Value: nil, 
						Scope: Scope, 
						Real: "e6", 
						Extern: extern, 
						ArgNum: nargs,
						ArgumentTypeManifest: ManifestEntries,
						BasinSize: BasinSize,
					}	
				} else {
					FuncVar = Variable_Static{
						Name: name, 
						Type: NUMBER16,
						Type2: rtype,
						Pointer: true,
						PointerLength: _ptrlen,
						Value: nil, 
						Scope: Scope, 
						Real: "e6", 
						Extern: extern, 
						ArgNum: nargs,
						ArgumentTypeManifest: ManifestEntries,
						PointingToStruct: PTS,
						BasinSize: BasinSize,
					}	
				}
				Variables = append(Variables, FuncVar)

				noreturn := false
				if peek(0).Value == "__attribute__" {
					attrs := _PARSE_ATTR(name)
					for _, attr := range attrs {
						switch attr {
						case "noreturn":
							noreturn = true
						default:
							error.Warning(34, "'" + attr + "' not allowed here", peek(-3), &tokens)
						}
					}
				}

				if peek(0).Type == shared.TokSemi {
					expect(shared.TokSemi)
					continue
				}
	
				expect(shared.TokLCurly)	

				var Children = []shared.Token {}
				ending := -1

				depth := 1
				for j := i; j < len(tokens); j++ {
					if tokens[j].Type == shared.TokRCurly {
						depth--
						if depth == 0 {
							ending = j
							break
						} else {
							Children = append(Children, tokens[j])
						}	
					} else if tokens[j].Type == shared.TokLCurly {
						depth++
						Children = append(Children, tokens[j])
					} else {
						Children = append(Children, tokens[j])
					}
				}
				if ending == -1 {
					error.Error(2, "'}'", tokens[i], &tokens)
				} else {
					i = ending
				}
			
				expect(shared.TokRCurly)
	
				Write(name + ":", false)	

				if len(Children) > 0 {
					level = 1
					if static == false {
						PreWrite(".global " + name, false)
					}
					topLevelName = name
					if name != "_start" && noreturn == false {
						Write("pop e11", true)	
					}
					if register > 0 {
						for r := nargs - 1; nargs > 0; nargs-- {
							Write("pop e" + fmt.Sprintf("%d", r), true)
							r--
						}
					}

					Write("push fp", true)

					FPtr := LookupVariableDirect(name, 1)
					(*FPtr).HasBasin = true
	
					Write("mov r12, _builtin_lcc_basin_" + name, true)
					Write("sub fp, fp, r12", true);

					if noreturn == false {
						Write("push e11", true)
					}

					for _, UnpackOrder := range UnpackOrders {
						__rn := UnpackOrder.Label
						__arg_reg := UnpackOrder.Register
						__ptr := UnpackOrder.Pointer
						__rtype := UnpackOrder.Type

						switch __ptr {
						case false:
							switch __rtype {
							case NUMBER8:
								Write("mov r1, " + __rn, true)
								Write("str r1, " + __arg_reg, true)
							case STRING:
								Write("mov r1, " + __rn, true)
								Write("str r1, " + __arg_reg, true)
							case NULL:
								Write("mov r1, " + __rn, true)
								Write("str_ptr r1, " + __arg_reg, true)	
							case NUMBER16:
								Write("mov r1, " + __rn, true)
								Write("str16 r1, " + __arg_reg, true)
							case NUMBER32:
								Write("mov r1, " + __rn, true)
								Write("str32 r1, " + __arg_reg, true)
							}
						case true:
							Write("mov r1, " + __rn, true)
							Write("str_ptr r1, " + __arg_reg, true)
						}	
					}

					_CURRENT_INFUNCTION = FuncVar	
					ParseExpyL1(Children, 0, fscope)

					if noreturn == false {
						Write("pop e11", true)
					}

					Write("pop fp", true)

					if noreturn == false {
						Write("ret", true)
					}
					IDCounter++
					topLevelName = ""
					level = 0
				}
			case shared.TokIdent:
				if peek(0).Value == "__attribute__" {
					attrs = _PARSE_ATTR(name)
					for _, attr := range attrs {
						switch attr {
						case "require_const":
							allow_nonconst = false	
						default:
							error.Warning(34, "'" + attr + "' not allowed here", peek(-3), &tokens)
						}
					}	
				} else {
					// Error out
					expect(shared.TokEOF)
				}
				fallthrough
			case shared.TokEqual:	
				expect(shared.TokEqual)	
				switch _type {
				case "void":
					error.Error(7, "'void'", _typetoken, &tokens)
				case "int":
					_i := 0

					rn := "fp + "

					switch IS_LOCAL {
					case false:
						rn = "var_" + fmt.Sprintf("%d", IDCounter)
					case true:
						switch ptr {
						case true:
							Displacement := ReturnDisplacement(NUMBER32)
							rn += fmt.Sprintf("%d", Displacement)
						case false:
							Displacement := ReturnDisplacement(rtype)
							rn += fmt.Sprintf("%d", Displacement)
						}	
					}

					if IS_LOCAL == true {
						FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
						switch ptr {
						case true:
							(*FPtr).BasinSize += 4
						case false:
							switch rtype {
							case NUMBER8:
								(*FPtr).BasinSize += 1
							case NUMBER16:
								(*FPtr).BasinSize += 2
							case NUMBER32, NULL, STRING:
								(*FPtr).BasinSize += 4
							}
						}
					}

					IDCounter++	

					if allow_nonconst == false {
						res := 0
						WritePre(rn + ":", false)
						res, _i = ParseNumberExpyDirect(tokens, i, Scope)

						if res == -1 {
							goto EQU_RTYPE_DONE
						}

						
						if ptr == false {
							switch rtype {
							case NUMBER8, STRING:
								r_res := uint8(res)
								if res > math.MaxUint8 || res < 0 {
									error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned short short int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(0), &tokens)
								}
								WritePre(".byte " + fmt.Sprintf("0x%02x", r_res), true)
							case NUMBER16:
								r_res := uint16(res)
								if res > math.MaxUint16 || res < 0 {
									error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned short int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(0), &tokens)
								}
								WritePre(".word " + fmt.Sprintf("0x%04x", r_res), true)
							case NUMBER32:
								r_res := uint32(res)
								if res > math.MaxUint32 || res < 0 {
									error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned long int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(0), &tokens)
								}
								WritePre(".dword " + fmt.Sprintf("0x%08x", r_res), true)
							}
						} else {
							WritePre(".ptr " + fmt.Sprintf("0x%x", res), true)
						}
						EQU_RTYPE_DONE:
					} else {
						if IS_LOCAL == false {
							WritePre(rn + ":", false)
						}
						var ATMEntry ArgumentTypeManifestEntry
						if ptr == false {
							ATMEntry.Type = rtype
							ATMEntry.Pointer = ptr
							ATMEntry.PointerLength = _ptrlen
						} else {
							ATMEntry.Type2 = rtype
							ATMEntry.Pointer = ptr
							ATMEntry.PointerLength = _ptrlen
						}
						_i = ParseExpy(tokens, i, Scope, "r4", ATMEntry)

						if ptr == false {
							switch rtype {
							case NUMBER8, STRING:
								if IS_LOCAL == false {
									WritePre(".byte 0x00", true)
								}
								Write("mov r7, " + rn, true)
								Write("str r7, r4", true)
							case NUMBER16:
								if IS_LOCAL == false {
									WritePre(".word 0x0000", true)
								}
								Write("mov r7, " + rn, true)
								Write("str16 r7, r4", true)
							case NUMBER32:
								if IS_LOCAL == false {
									WritePre(".dword 0x00000000", true)
								}
								Write("mov r7, " + rn, true)
								Write("str32 r7, r4", true)
							}
						} else {
							if IS_LOCAL == false {
								WritePre(".ptr 0x00", true)
							}
							Write("mov r7, " + rn, true)
							Write("str_ptr r7, r4", true)
						}
					}
					i = _i


					var val any
					if ptr == true {	
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: rtype, Value: val, Pointer: true, PointerLength: _ptrlen, Real: rn, Scope: Scope, Const: constant})
					} else {
						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: val, Pointer: false, Real: rn, Scope: Scope, Const: constant})
					}
				case "char":
					// TODO: add inetgers to this as well
					str, end := StringParse(tokens, i)	
					if ptr == true {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						rn2 := ""

						WritePre(rn + ":", false)
						WritePre(".asciz \"" + str + "\"", true)

						if IS_LOCAL == true {
							FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
							(*FPtr).BasinSize += 4
							rn2 = "fp + " + fmt.Sprintf("%d", ReturnDisplacement(NUMBER32))	

							Write("mov r7, " + rn2, true)
							Write("mov r4, " + rn, true)
							Write("str_ptr r7, r4", true)
						} else {
							rn2 = "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn2 + ":", false)
							WritePre(".ptr " + rn, true)
						}
	
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Value: str, Type2: STRING, Pointer: true, PointerLength: _ptrlen, Real: rn2, Scope: Scope, Const: constant})	
					} else {	
						if len(str) > 1 {
							error.Error(5, "'char' with an expression of type 'char*'", tokens[i], &tokens)
						}

						rn := ""
						if IS_LOCAL == true {
							FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
							(*FPtr).BasinSize += 1
							rn = "fp + " + fmt.Sprintf("%d", ReturnDisplacement(NUMBER8))
						} else {
							rn = "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn + ":", false)
							WritePre(".byte " + fmt.Sprintf("0x%02x", str[0]), true)
						}

						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: str, Pointer: false, Scope: Scope, Const: constant, Real: rn})
					}
					i = end + 1
				default:
					var TypeEntry TypeMapEntry
					found := false
					for _, Type := range TypeMap {
						if Type.Name == _type {
							found = true
							TypeEntry = Type
							break
						}
					}
					if found == false {
						error.Error(46, "'" + _type + "'", _type_tok, &tokens)
						for j := i; j < len(tokens); j++ {
							if tokens[j].Type == shared.TokSemi {
								i = j + 1
								break
							}
						}
						continue
					}
					rtype = TypeEntry.ReferralType

					switch TypeEntry.ReferralType {
					case STRUCT:
						if ptr == false {
							error.UnimplementedMessage("direct assignment of structs is not supported.")
						} else {
							_i := 0;
							rn := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++

							var ATMEntry ArgumentTypeManifestEntry
							if ptr == false {
								ATMEntry.Type = rtype
								ATMEntry.Pointer = ptr
								ATMEntry.PointerLength = _ptrlen
							} else {
								ATMEntry.Type2 = rtype
								ATMEntry.Pointer = ptr
								ATMEntry.PointerLength = _ptrlen
							}

							if IS_LOCAL == false {
								WritePre(rn + ":", false)
								WritePre(".ptr 0x00", true)
							} else {
								rn = "fp + " + fmt.Sprintf("%d", ReturnDisplacement(NUMBER32))
								FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
								(*FPtr).BasinSize += 4	
							}
								
							// TODO: var lookup for toplevel
							_i = ParseExpy(tokens, i, Scope, "r4", ATMEntry)
	
							Write("mov r7, " + rn, true)
							Write("str_ptr r7, r4", true)

							i = _i

							var val any
							Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: rtype, Value: val, Pointer: true, PointerLength: _ptrlen, Real: rn, Scope: Scope, PointingToStruct: _type, Const: constant})
						}
					case NUMBER8, NUMBER16, NUMBER32, STRING, NULL:
						_i := 0
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++	

						if allow_nonconst == false {
							res := 0
							WritePre(rn + ":", false)
							res, _i = ParseNumberExpyDirect(tokens, i, Scope)

							if res == -1 {
								goto EQU_RTYPE_DONE2
							}

							switch rtype {
							case NUMBER8, STRING:
								r_res := uint8(res)
								if res > math.MaxUint8 || res < 0 {
									error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned short short int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(-1), &tokens)
								}
								WritePre(".byte " + fmt.Sprintf("0x%02x", r_res), true)
							case NUMBER16:
								r_res := uint16(res)
								if res > math.MaxUint16 || res < 0 {
									error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned short int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(-1), &tokens)
								}
								WritePre(".word " + fmt.Sprintf("0x%04x", r_res), true)
							case NUMBER32:
								r_res := uint32(res)
								if res > math.MaxUint32 || res < 0 {
									error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned long int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(-1), &tokens)
								}
								WritePre(".dword " + fmt.Sprintf("0x%08x", r_res), true)
							}	
							EQU_RTYPE_DONE2:
						} else {	
							var ATMEntry ArgumentTypeManifestEntry
							if ptr == false {
								ATMEntry.Type = rtype
								ATMEntry.Pointer = ptr
								ATMEntry.PointerLength = _ptrlen
							} else {
								ATMEntry.Type2 = rtype
								ATMEntry.Pointer = ptr
								ATMEntry.PointerLength = _ptrlen
							}

							_i = ParseExpy(tokens, i, Scope, "r4", ATMEntry)

							if IS_LOCAL == true {
								FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
								Displacement := uint32(0)
								if ptr == false {
									Displacement = ReturnDisplacement(TypeEntry.ReferralType)
									rn = "fp + " + fmt.Sprintf("%d", Displacement)
									switch TypeEntry.ReferralType {
									case NUMBER8, STRING, NULL:
										(*FPtr).BasinSize += 1
									case NUMBER16:
										(*FPtr).BasinSize += 2
									case NUMBER32:
										(*FPtr).BasinSize += 4
									}
								} else {
									Displacement = ReturnDisplacement(NUMBER32)
									rn = "fp + " + fmt.Sprintf("%d", Displacement)
									(*FPtr).BasinSize += 4
								}

								if ptr == false {
									switch rtype {
									case NUMBER8, STRING, NULL:
										Write("mov r7, " + rn, true)
										Write("str r7, r4", true)
									case NUMBER16:
										Write("mov r7, " + rn, true)
										Write("str16 r7, r4", true)
									case NUMBER32:
										Write("mov r7, " + rn, true)
										Write("str32 r7, r4", true)
									}
								} else {
									Write("mov r7, " + rn, true)
									Write("str_ptr r7, r4", true)
								}
							} else {
								// TODO: fix this
								WritePre(rn + ":", false)
								if ptr == false {
									switch rtype {
									case NUMBER8, STRING:
										WritePre(".byte 0x00", true)
										Write("mov r7, " + rn, true)
										Write("str r7, r4", true)
									case NUMBER16:
										WritePre(".word 0x0000", true)
										Write("mov r7, " + rn, true)
										Write("str16 r7, r4", true)
									case NUMBER32:
										WritePre(".dword 0x00000000", true)
										Write("mov r7, " + rn, true)
										Write("str32 r7, r4", true)
									}
								} else {
									WritePre(".ptr 0x00", true)
									Write("mov r7, " + rn, true)
									Write("str_ptr r7, r4", true)
								}
							}	
						}

						i = _i

						var val any
						if ptr == true || TypeEntry.Pointer == true {	
							Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: rtype, Value: val, Pointer: true, PointerLength: _ptrlen + TypeEntry.PointerLength, Real: rn, Scope: Scope, Const: constant})
						} else {
							Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: val, Pointer: false, Real: rn, Scope: Scope, Const: constant})
						}
					}	
				}

				expect(shared.TokSemi)
			case shared.TokSemi:
				expect(shared.TokSemi)

				switch _type {
				case "int":
					if ptr == true {
						rn := ""
						if IS_LOCAL == true {
							// TODO: differentiate between 16 and 32 bit pointers to save space
							rn = "fp + " + fmt.Sprintf("%d", ReturnDisplacement(NUMBER32))
							FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
							(*FPtr).BasinSize += 4
							// Uninitialized, so no store
						} else {
							rn = "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn + ":", false)
							WritePre(".ptr 0x00000000", true)	
						}	

						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: 0, Pointer: true, PointerLength: _ptrlen, Real: rn, Scope: Scope, Const: constant})
					} else {
						rn := ""
						if IS_LOCAL == true {
							// TODO: differentiate between 16 and 32 bit pointers to save space
							rn = "fp + " + fmt.Sprintf("%d", ReturnDisplacement(NUMBER32))
							FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
							(*FPtr).BasinSize += 4
							// Uninitialized, so no store
						} else {
							rn = "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn + ":", false)

							switch rtype {
							case NUMBER8, STRING:
								WritePre(".byte 0x00", true)
							case NUMBER16:
								WritePre(".word 0x0000", true)
							case NUMBER32:
								WritePre(".dword 0x00000000", true)
							}
						}

						Variables = append(Variables, Variable_Static{Name: name, Real: rn, Type: rtype, Value: 0, Pointer: false, Scope: Scope, Const: constant})	
					}
				case "char":
					if ptr == true {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						rn2 := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: STRING, Value: "", Pointer: true, PointerLength: _ptrlen, Real: rn2, Scope: Scope, Const: constant})
						WritePre(rn + ":", false)
						WritePre(".asciz \"\"", true)
						WritePre(rn2 + ":", false)
						WritePre(".ptr " + rn, true)
					} else {	
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: "", Pointer: false, Scope: Scope, Const: constant})
					}
				case "void":
					if ptr == true {
						if extern == false {
							rn := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							rn2 := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn2 + ":", false)
							WritePre(".ptr " + rn, true)
							WritePre(rn + ":", false)
							Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: NULL, Value: nil, Pointer: true, PointerLength: _ptrlen, Real: rn, Scope: Scope, Const: constant})
						} else {
							rn := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn + ":", false)
							WritePre(".ptr " + name, true)
							Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: NULL, Value: nil, Pointer: true, Real: rn, PointerLength: _ptrlen, Scope: Scope, Const: constant})	
						}
					}
				default:
					var TypeEntry TypeMapEntry
					found := false
					for _, Type := range TypeMap {
						if Type.Name == _type {
							found = true
							TypeEntry = Type
							break
						}
					}
					if found == false {
						error.Error(46, "'" + _type + "'", _type_tok, &tokens)
						for j := i; j < len(tokens); j++ {
							if tokens[j].Type == shared.TokSemi {
								i = j + 1
								break
							}
						}
						continue
					}
					rtype = TypeEntry.ReferralType

					switch TypeEntry.ReferralType {
					case STRUCT:
						if ptr == false {
							rn := ""
							if IS_LOCAL == true {
								rn = "fp + " + fmt.Sprintf("%d", CurrentDisplacement)
								CurrentDisplacement += uint32(TypeEntry.EmbeddedStruct.StructTotalSize)
								FPtr := LookupVariableDirect(_CURRENT_INFUNCTION.Name, 1)
								(*FPtr).BasinSize += uint32(TypeEntry.EmbeddedStruct.StructTotalSize)
								
							} else {
								rn = "var_" + fmt.Sprintf("%d", IDCounter)
								IDCounter++

								WritePre(rn + ":", false)
								WritePre(".pad " + fmt.Sprintf("%d", TypeEntry.EmbeddedStruct.StructTotalSize), true)
							}

							StructVar := TypeEntry.EmbeddedStruct

							StructVar.Name = name
							StructVar.Real = rn
							StructVar.Type = STRUCT
							StructVar.Scope = Scope
							StructVar.StructTotalSize = TypeEntry.EmbeddedStruct.StructTotalSize

							Variables = append(Variables, StructVar)
						} else {
							rn := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++

							WritePre(rn + ":", false)
							WritePre(".ptr 0", true)

							fmt.Println(_type)

							Variables = append(Variables, Variable_Static {
								Name: name,
								Real: rn,
								Type: NUMBER16,
								Pointer: true,
								PointerLength: _ptrlen,
								Type2: rtype,
								Scope: Scope,
								Const: constant,
								PointingToStruct: _type,
							})
						}	
					}
				}
			case shared.TokLBracket:
				expect(shared.TokLBracket)
				// Length next
				if peek(0).Type == shared.TokIdent {
					i += 2
					error.UnimplementedMessage("variable-length arrays are not supported.")
				}
				length := expect(shared.TokNumber)
				length_real, _ := strconv.ParseInt(length, 0, 64)
				expect(shared.TokRBracket)

				rn := "var_" + fmt.Sprintf("%d", IDCounter)
				IDCounter++
				Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: nil, Pointer: false, Real: rn, Scope: Scope, Const: constant, ArgNum: int(length_real) })
				WritePre(rn + ":", false)

				if rtype == NUMBER8 || rtype == STRING {
					length_real = length_real
				} else if rtype == NUMBER16 {
					length_real = length_real * 2
				} else if rtype == NUMBER32 {
					length_real = length_real * 4
				}

				WritePre(".pad " + fmt.Sprintf("%d", length_real), true)

				expect(shared.TokSemi)
			default:
				error.Error(1, "'" + peek(0).Value + "'", _typetoken, &tokens)
			}	
		}
	}
}
