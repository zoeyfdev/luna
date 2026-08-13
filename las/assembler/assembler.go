package assembler

import (
	"strings"
	"strconv"
	"las/error"
	"path/filepath"
	"os"
)

var Bits32 bool = false
var ForcedSize int64 = 0
var PIE bool

var DataBuffer []byte
var TextBuffer []byte
var ExtendedDataBuffer []byte

func write(b []byte) {
	// sections aren't really a thing anymore
	TextBuffer = append(TextBuffer, b...)	
}

func isRegister(word string) byte {
	switch word {
	case "r0":
		return 0x00
	case "r1":
		return 0x01
	case "r2":
		return 0x02
	case "r3":
		return 0x03
	case "r4":
		return 0x04
	case "r5":
		return 0x05
	case "r6":
		return 0x06
	case "r7":
		return 0x07
	case "r8":
		return 0x08
	case "r9":
		return 0x09
	case "r10":
		return 0x0a
	case "r11":
		return 0x0b
	case "r12":
		return 0x0c
	case "e0":
		return 0x0d
	case "e1":
		return 0x0e
	case "e2":
		return 0x0f
	case "e3":
		return 0x10
	case "e4":
		return 0x11
	case "e5":
		return 0x12
	case "e6":
		return 0x13
	case "e7":
		return 0x14
	case "e8":
		return 0x15
	case "e9":
		return 0x16
	case "e10":
		return 0x17
	case "e11":
		return 0x18
	case "e12":
		return 0x19
	case "e13":
		return 0x1a
	case "e14":
		return 0x1b
	case "sp":
		return 0x1c
	case "pc":
		return 0x1d
	case "irv":
		return 0x1e
	case "ir":
		return 0x1f
	case "b":
		return 0x20
	case "fp":
		return 0x21
	default:
		return 0xff
	}
}

func parse(text string) ([]byte, bool) {
	// Check for number
	if _, err := strconv.ParseInt(text, 0, 64); err == nil {
		num, _ := strconv.ParseInt(text, 0, 64)
		if Bits32 == false {
			H := byte(num >> 8)
			L := byte(num & 0xFF)
			return []byte{H, L}, false
		} else {
			HH := byte(num >> 24)
			HL := byte(num >> 16)
			LH := byte(num >> 8)
			LL := byte(num & 0xFF)
			return []byte{HH, HL, LH, LL}, false	
		}
	}
	if isRegister(text) != 0xff {	
		return []byte{byte(isRegister(text))}, false
	}
	if string(text[0]) == "\"" {
		if string(text[len(text) - 1]) != "\"" {
			error.Error(6, "")
		}

		text = strings.Trim(text, "\"")

		text = strings.ReplaceAll(text, "\\0", "\000")
		text = strings.ReplaceAll(text, "\\n", "\n")
		text = strings.ReplaceAll(text, "\\r", "\r")
		if Bits32 == false {
			if len(text) > 2 {
				error.Error(5, "'" + text + "'")
			} else if len(text) == 1 {
				text = string(byte(00)) + text
			}
		} else {
			if len(text) > 4 {
				error.Error(5, "'" + text + "'")	
			} else if len(text) == 1 {
				text = string(byte(00)) + string(byte(00)) + string(byte(00)) + text
			} else if len(text) == 2 {
				text = string(byte(00)) + string(byte(00)) + text
			} else if len(text) == 3 {
				text = string(byte(00)) + text
			}
		}	
		return []byte(text), false
	}	
	return append([]byte("LR_"+text), 0x00), true
}

func formatString(text string) string {
	var replace = [][2]string {
		{"\\0", "\000"},
		{"\\n", "\n"},
		{"\\r", "\r"},
		{"\\033", "\033"},
	}
	for _, pair := range replace {
		text = strings.ReplaceAll(text, pair[0], pair[1])
	}
	return text
}


func Lex(text string) []Token {
   	var Tokens []Token
	var Buffer string
	var Str bool
	var Line int = 1
	var InComment bool = false

	Add := func(Text string) {
		if Text == "" { return }
		Tokens = append(Tokens, Token{
			Value: Text,
			Line: Line,
		})
	}

	for i := 0; i < len(text); i++ {
		current := text[i]

		switch current {
		case '"':
			if InComment == false {
				Buffer += "\""
				if Str == false {
					Str = true
				} else {
					Str = false
					Add(Buffer)
					Buffer = ""
				}
			}
		case 0x0A:
			if InComment == true { InComment = false }
			if Str == false && InComment == false {
				Add(Buffer)
				Buffer = "" 
			}
			Add("\n")
			Line++
		case ' ':
			if Str == false && InComment == false {
				Add(Buffer)
				Buffer = ""
			} else {
				if Str == true && InComment == false {
					Buffer += string(current)
				}
			}
		case '\t':
			if Str == false && InComment == false {
				Add(Buffer)
				Buffer = ""
			}	
		case '/':
			if text[i + 1] == '/' {
				i++
				InComment = true
			} else {
				if InComment == false {
					Buffer += string(current)
				}
			}
		case ';':
			InComment = true
		case '#':
			if text[i + 1] == ' ' {
				i++
				InComment = true
			} else {
				if InComment == false {
					Buffer += string(current)
				}
			}
		default:
			if InComment == false {
				Buffer += string(current)
			}
		}
	}

	if Buffer != "" && InComment == false {
		Add(Buffer)
	}
	
	return Tokens
}

var non_organic bool
var PJL bool
var clabel string

func Assemble(text string) {
	Tokens := Lex(text)
	words := []string {}

	for _, T := range Tokens {
		words = append(words, T.Value) // because I can't be bothered to change out words[x] to Token.Value
	}

	for i := 0; i < len(words); i++ {
		words[i] = strings.TrimSuffix(words[i], ",")
	}

	for i := 0; i < len(words); i++ {
		switch words[i] {
		case "#define":
			alias := words[i + 1]
			actual := words[i + 2]
			words = append(words[:i], words[i + 3:]...)
			for j := 0; j < len(words); j++ {
				if words[j] == alias {
					words[j] = actual
				}
			}
		case "#pragma":
			switch words[i + 1] {
			case "size":
				size, err := strconv.ParseInt(words[i + 2], 0, 64)
				if err != nil {
					error.Error(11, "")
					break
				}
				ForcedSize = size
				i++
			default:
				error.Warning(12, "'" + words[i + 1] + "'")	
			}
			i++
		}
	}

	for i := 0; i < len(words); i++ {
		if strings.HasSuffix(words[i], ":") && !strings.Contains(words[i], "\"") {
			end := len(words)
			for j := i + 1; j < len(words); j++ {
				if strings.HasSuffix(words[j], ":") {
					end = j
					break
				}
			}

			words[i] = strings.ReplaceAll(words[i], ":", "")
			if Bits32 == false {
				write(append([]byte("LD16_" + words[i]), 0x00))
			} else {
				write(append([]byte("LD32_" + words[i]), 0x00))
			}
			clabel = words[i]
			tocompile := words[i + 1:end]
			if len(tocompile) > 0 {
				Assemble(strings.Join(tocompile, " "))
			}

			i = end - 1
			continue
		}

		words[i] = strings.ToLower(words[i])


		switch words[i] {
		case "\n":
			error.Line++
			continue
		case "mov":
			write([]byte{0x01})

			var mode byte
			if isRegister(words[i + 2]) == 0xff {
				mode = 0x01
			} else {
				if words[i + 3] == "+" || words[i + 3] == "-" {
					mode = 0x03
				} else {
					mode = 0x02
				}
			}
			write([]byte{mode})

			dst := isRegister(words[i+1])
			if dst == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			write([]byte{dst})

			if mode == 0x02 {
				src := isRegister(words[i + 2])
				write([]byte{src})
			} else if mode == 0x03 {
				src := isRegister(words[i + 2])
				write([]byte{src})

				_type := words[i + 3]
				if _type == "+" {
					write([]byte{ 0x01 })
				} else if _type == "-" {
					write([]byte{ 0x02 })
				}

				// TODO: add support for PIEs
				value, _ := parse(words[i + 4])
				write(value)
				
				i += 2
			} else {
				value, label := parse(words[i + 2])
	
				write(value)
				if PIE == true && label == true {	
					if non_organic == true {
						write([]byte {0x0d, 0x1a, dst, 0x1b})
					} else {
						write([]byte {0x0d, dst, dst, 0x1b})
					}
				}
			}
			i += 2
		case "hlt":
			write([]byte{0x02})
		case "jmp":
			_, label := parse(words[i + 1])
			if PIE == false || label == false {
				write([]byte{0x03})

				if isRegister(words[i + 1]) == 0xff {
					write([]byte{0x01})
				} else {
					write([]byte{0x02})
				}

				value, _ := parse(words[i + 1])
				write(value)
			} else {
				non_organic = true
				Assemble(`mov e13, ` + words[i + 1])
				non_organic = false
				Assemble(`jmp e13`)
			}
			i = i + 1
		case "int":
			write([]byte{0x04})
			value, _ := parse(words[i+1])	
			if Bits32 == false {
				if len(value) > 2 {
					error.Error(3, "'" + string(value) + "'")
				}
			} else {
				if len(value) > 4 {
					error.Error(3, "'" + string(value) + "'")
				}
			}
			write(value)
			i = i + 1
		case "jnz":
			_, label := parse(words[i + 2])
			if PIE == false || label == false {
				write([]byte{0x05})

				if isRegister(words[i+2]) == 0xff {
					write([]byte{0x01})
				} else {
					write([]byte{0x02})
				}

				register := isRegister(words[i+1])
				if register == 0xff {
					error.Error(2, "'"+words[i+1]+"'")
				}
				write([]byte{register})

				value, _ := parse(words[i+2])
				write(value)
			}  else {
				register := isRegister(words[i+1])
				if register == 0xff {
					error.Error(2, "'"+words[i+1]+"'")
				}
				non_organic = true
				Assemble(`mov e13, ` + words[i + 2])
				non_organic = false
				Assemble(`jnz ` + words[i + 1] + `, e13`)
			}
			i = i + 2
		case "nop":
			write([]byte{0x06})
		case "cmp":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x07})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "jz":
			_, label := parse(words[i + 2])
			if PIE == false || label == false {
				write([]byte{0x08})

				if isRegister(words[i+2]) == 0xff {
					write([]byte{0x01})
				} else {
					write([]byte{0x02})
				}

				register := isRegister(words[i+1])
				if register == 0xff {
					error.Error(2, "'"+words[i+1]+"'")
				}
				write([]byte{register})

				value, _ := parse(words[i+2])
				write(value)
			} else {
				register := isRegister(words[i+1])
				if register == 0xff {
					error.Error(2, "'"+words[i+1]+"'")
				}
				non_organic = true
				Assemble(`mov e13, ` + words[i + 2])
				non_organic = false
				Assemble(`jz ` + words[i + 1] + `, e13`)
			}
			i = i + 2
		case "inc":
			write([]byte{0x09})
			reg := isRegister(words[i+1])
			if reg == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			write([]byte{reg})
			i = i + 1
		case "dec":
			write([]byte{0x0a})
			reg := isRegister(words[i+1])
			if reg == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			write([]byte{reg})
			i = i + 1
		case "push":
			_, label := parse(words[i + 1])
			if PIE == false || label == false {
				write([]byte{0x0b})
				if isRegister(words[i+1]) == 0xff {
					write([]byte{0x01})
				} else {
					write([]byte{0x02})
				}
				value, _ := parse(words[i + 1])
				write(value)
			} else {
				non_organic = true
				Assemble(`mov e13, ` + words[i + 1])
				non_organic = false
				Assemble(`push e13`)
			}
			i = i + 1
		case "pop":
			write([]byte{0x0c})
			reg := isRegister(words[i+1])
			if reg == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			write([]byte{reg})
			i = i + 1
		case "add":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x0d})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "sub":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x0e})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "mul":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x0f})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "div":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x10})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "igt":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x11})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "ilt":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x12})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "and":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x13})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "or":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x14})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "not":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			write([]byte{0x15})
			write([]byte{check})
			write([]byte{one})
			i = i + 2
		case "xor":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x16})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "lod", "lode":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if words[i] == "lode" {
				Assemble("add e13, " + words[i + 1] + ", e14")
				write([]byte{0x17})
				write([]byte{0x1a})
				write([]byte{one})
			} else {
				write([]byte{0x17})
				write([]byte{check})
				write([]byte{one})
			}
			i = i + 2
		case "str16", "str16e": // strfe
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if words[i] == "str16e" {
				Assemble("add e13, " + words[i + 1] + ", e14")
				write([]byte{0x18})
				write([]byte{0x1a})
				write([]byte{one})
			} else {
				write([]byte{0x18})
				write([]byte{check})
				write([]byte{one})
			}
			i = i + 2
		case "lod16", "lod16e": // lodfe
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if words[i] == "lod16e" {
				Assemble("add e13, " + words[i + 1] + ", e14")
				write([]byte{0x19})
				write([]byte{0x1a})
				write([]byte{one})
			} else {
				write([]byte{0x19})
				write([]byte{check})
				write([]byte{one})
			}
			i = i + 2
		case "set":
			mode := words[i + 1]

			switch mode {
			case "16":
				write([]byte{0x1a, 0x00})	
			case "32":
				write([]byte{0x1a, 0x01})	
			}
			i++
		case "str", "stre":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if words[i] == "stre" {
				Assemble("add e13, " + words[i + 1] + ", e14")
				write([]byte{0x1b})
				write([]byte{0x1a})
				write([]byte{one})
			} else {
				write([]byte{0x1b})
				write([]byte{check})
				write([]byte{one})
			}
			i = i + 2
		case "shl":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x1c})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "shr":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x1d})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "str32", "str32e": // strfe
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if words[i] == "str32e" {
				Assemble("add e13, " + words[i + 1] + ", e14")
				write([]byte{0x1f})
				write([]byte{0x1a})
				write([]byte{one})
			} else {
				write([]byte{0x1f})
				write([]byte{check})
				write([]byte{one})
			}
			i = i + 2
		case "lod32", "lod32e": // lodfe
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if words[i] == "lod32e" {
				Assemble("add e13, " + words[i + 1] + ", e14")
				write([]byte{0x1e})
				write([]byte{0x1a})
				write([]byte{one})
			} else {
				write([]byte{0x1e})
				write([]byte{check})
				write([]byte{one})
			}
			i = i + 2
		case "mod":
			check := isRegister(words[i+1])
			one := isRegister(words[i+2])
			two := isRegister(words[i+3])
			if check == 0xff {
				error.Error(2, "'"+words[i+1]+"'")
			}
			if one == 0xff {
				error.Error(2, "'"+words[i+2]+"'")
			}
			if two == 0xff {
				error.Error(2, "'"+words[i+3]+"'")
			}
			write([]byte{0x20})
			write([]byte{check})
			write([]byte{one})
			write([]byte{two})
			i = i + 3
		case "lodf", "strf":
			error.Warning(13, "'" + words[i] + "'")	
			if Bits32 == false {
				Assemble(words[i][0:3] + "16 " + words[i + 1] + " " + words[i + 2])
			} else {
				Assemble(words[i][0:3] + "32 " + words[i + 1] + " " + words[i + 2])
			}
			
			i += 2
		case "lod_ptr", "str_ptr":			
			if Bits32 == false {
				Assemble(words[i][0:3] + "16 " + words[i + 1] + " " + words[i + 2])
			} else {
				Assemble(words[i][0:3] + "32 " + words[i + 1] + " " + words[i + 2])
			}
			
			i += 2
		case "call":
			label := words[i + 1]
			if PIE == false {
				if Bits32 == false {
					Assemble(`
					mov e11, pc                  
					mov r0, 20                   
					add e11, e11, r0             
					push e11                     
					jmp	` + label)             
				} else {
					Assemble(`
					mov e11, pc
					mov r0, 24
					add e11, e11, r0
					push e11
					jmp	` + label) 
					// 4 bytes
					// 7 bytes
					// 4 bytes
					// 3 bytes
					// 6 bytes << 14 bytes
				}
			} else {
				if Bits32 == false {
					Assemble(`
					mov e11, pc                  
					mov r0, 20                   
					add e11, e11, r0             
					push e11                     
					jmp	` + label)             
				} else {
					Assemble(`
					mov e11, pc
					mov r0, 32
					add e11, e11, r0
					push e11
					jmp	` + label) 
					// 4 bytes
					// 7 bytes
					// 4 bytes
					// 3 bytes
					// 6 bytes << 14 bytes
				}
			}
			
			i = i + 1
		case "ret":
			if (clabel == "_start" || clabel == "") && PIE == true {
				Assemble(`
				mov r1, 1
				int 0x04
				`)	
			} else {
				Assemble(`jmp e11`)
			}
		case "pusha":
			Assemble(`
				push r0
				push r1
				push r2
				push r3
				push r4
				push r5
				push r6
				push r7
				push r8
				push r9
				push r10
				push r11
				push r12
				push e0
				push e1
				push e2
				push e3
				push e4
				push e5
				push e6
				push e7
				push e8
				push e9
				push e10
				push e11
				push e12
			`)
		case "popa":
			Assemble(`
				pop e12
				pop e11
				pop e10
				pop e9
				pop e8
				pop e7
				pop e6
				pop e5
				pop e4
				pop e3
				pop e2
				pop e1
				pop e0
				pop r12
				pop r11
				pop r10
				pop r9
				pop r8
				pop r7
				pop r6
				pop r5
				pop r4
				pop r3
				pop r2
				pop r1
				pop r0
			`)
		case "syscall":
			Assemble(`int 0x4`)
		case ".ascii":	
			var value string	
			var tokens = []string {}
			
			if string(words[i+1][0]) != "\"" {
				error.Error(7, "'" + words[i+1] + "'")
			}
			if strings.HasSuffix(words[i + 1], "\"") {
				value = strings.Trim(words[i + 1], "\"")
				value = formatString(value)
				write([]byte(value))
				i = i + 1
				continue
			}
			
			ending := 0
			for j := i + 1; j < len(words); j++ {
				tokens = append(tokens, words[j])
				if strings.HasSuffix(words[j], "\"") {
					ending = j
					break
				}
			}
			if ending == 0 {
				error.Error(6, "'" + words[i + 1] + "'")
			}
			
			tokens[0] = strings.TrimPrefix(tokens[0], "\"")
			tokens[len(tokens) - 1] = strings.TrimSuffix(tokens[len(tokens) - 1], "\"")
			value = strings.Join(tokens, " ")
			value = formatString(value)
			write([]byte(value))
			i = ending
		case ".asciz":	
			var value string	
			var tokens = []string {}	

			if string(words[i+1][0]) != "\"" {
				error.Error(7, "'" + words[i+1] + "'")
			}
			if strings.HasSuffix(words[i + 1], "\"") {
				value = strings.Trim(words[i + 1], "\"")
				value = formatString(value)
				value = value + string("\000")
				write([]byte(value))
				i = i + 1
				continue
			}
			
			ending := 0
			for j := i + 1; j < len(words); j++ {
				tokens = append(tokens, words[j])
				if strings.HasSuffix(words[j], "\"") {
					ending = j
					break
				}
			}
			if ending == 0 {
				error.Error(6, "'" + words[i + 1] + "'")
			}
			
			tokens[0] = strings.TrimPrefix(tokens[0], "\"")
			tokens[len(tokens) - 1] = strings.TrimSuffix(tokens[len(tokens) - 1], "\"")
			value = strings.Join(tokens, " ")
			value = formatString(value)
			value = value + string("\000")
			write([]byte(value))
			i = ending
		case ".bits":
			switch words[i + 1] {
			case "16":
				Bits32 = false
				write([]byte("L_16BIT"))
			case "32":
				Bits32 = true
				write([]byte("L_32BIT"))
			default:
				error.Error(9, "")
			}
			i++
		case ".embed":
			absSource, _ := filepath.Abs(error.File)
			base := filepath.Dir(absSource)
			raw := strings.ReplaceAll(words[i + 1], "\"", "")

			path := raw
			if !filepath.IsAbs(raw) {
				path = filepath.Join(base, raw)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				error.Error(1, "'" + path + "'")
				continue
			}
			write(data)
			i++
		case ".word":
			word := words[i + 1]
			num, err := strconv.ParseInt(word, 0, 64)
			if err != nil {
				error.Error(11, ", got '" + word + "'")
			}
			H := byte(num >> 8)
			L := byte(num & 0xFF)
			write([]byte{H, L})
			i++
		case ".dword":
			word := words[i + 1]
			num, err := strconv.ParseInt(word, 0, 64)
			if err != nil {
				error.Error(11, ", got '" + word + "'")
			}
			HH := byte(num >> 24)
			HL := byte(num >> 16)
			LH := byte(num >> 8)
			LL := byte(num & 0xFF)
			write([]byte{HH, HL, LH, LL})
			i++
		case ".fill":
			word := words[i + 1]
			num, err := strconv.ParseInt(word, 0, 64)
			if err != nil {
				error.Error(11, ", got '" + word + "'")
			}
			H := byte(num >> 8)
			L := byte(num & 0xFF)
			write(append([]byte("LP_"), []byte{H, L}...))
			i++
		case ".org":
			word := words[i + 1]
			num, err := strconv.ParseInt(word, 0, 64)
			if err != nil {
				error.Error(11, ", got '" + word + "'")
			}
			H := byte(num >> 8)
			L := byte(num & 0xFF)
			write(append([]byte("LO_"), []byte{H, L}...))
			i++
		case ".ptr":
			value, _ := parse(words[i + 1])
			write(value)
			i++
		case ".pad":
			word := words[i + 1]
			num, err := strconv.ParseInt(word, 0, 64)
			if err != nil {
				error.Error(11, ", got '" + word + "'")
			}
			for j := 0; int64(j) < num; j++ {
				write([]byte{0x00})
			}
			i++
		case ".global":
			label := words[i + 1]
			write([]byte("L_GLOBL_" + label))
			write([]byte{0x00})
			i++
		case ".byte":
			for j := i + 1; j < len(words); j++ {
				if words[j] != "\n" {
					num, err := strconv.ParseUint(strings.ReplaceAll(words[j], ",", ""), 0, 8)
					if err != nil {
						error.Error(11, ", got '" + words[j] + "'")
					}
					write([]byte{byte(num)})
				} else {
					i = j
					break
				}
			}
		case ".fhex":
			i++
			if len(words[i]) != 6 && len(words[i]) != 8 {
				error.Error(15, "'" + words[i] + "'")
				continue
			}
	
			if len(words[i]) == 8 || words[i][0:2] == "0x" {
				words[i] = words[i][2:]
			} else {
				error.Error(15, "'" + words[i] + "'")
				continue
			}	

			r, e1 := strconv.ParseInt(words[i][0:2], 16, 64)
			g, e2 := strconv.ParseInt(words[i][2:4], 16, 64)
			b, e3 := strconv.ParseInt(words[i][4:6], 16, 64)

			if e1 != nil || e2 != nil || e3 != nil {
				error.Error(15, "'" + words[i] + "'")
				continue
			}
			
			write([]byte{byte((r & 0xE0) | ((g >> 3) & 0x1C) | (b >> 6))})
		default:
			error.Error(4, "'" + words[i] + "'")
		}
	}
}
