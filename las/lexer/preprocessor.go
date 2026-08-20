package lexer

import (
	"las/error"
	"las/shared"
)

type DefineEntry struct {
	Name string
	ReplacementList []shared.Token
}

var Defines = []DefineEntry {
	DefineEntry { 
		Name: "__LCC__", 
		ReplacementList: []shared.Token { 
			shared.Token { 
				Value: "1", 
				File: "__DEFAULT__",
			}, 
		},
	},
}

func CheckDefined(Name string) (DefineEntry, bool) {
	for _, Entry := range Defines {
		if Entry.Name == Name {
			return Entry, true
		}
	}
	return DefineEntry{Name: "__NOTFOUND__"}, false
}

func Preprocessor(data string, predefs []DefineEntry) []shared.Token {
	out := []shared.Token {}

	for _, Predef := range predefs {
		Defines = append(Defines, Predef)	
	}

	again := false

	tokens := Tokenize(data, shared.File)

PREPROCESSOR_TOP:
	for i := 0; i < len(tokens); i++ {
		Token := tokens[i]

		switch Token.Value {
		case "#define":
			var Define DefineEntry
			i++

			Define.Name = tokens[i].Value
			i++

			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					i = j
					break
				}
				Define.ReplacementList = append(Define.ReplacementList, t)
			}
			Defines = append(Defines, Define)
			again = true
		case "#ifdef", "#ifndef":
			var FakeStream []shared.Token
			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					break
				}
				FakeStream = append(FakeStream, shared.Token{
					Value: t.Value,
					Line: t.Line,
					File: t.File,
				})
			}
			origin := FakeStream[0]

			dir := tokens[i].Value
			i++
			look_for := tokens[i].Value
			_, Found := CheckDefined(look_for)
			i++
			if (Found == true && dir == "#ifdef") || (Found == false && dir == "#ifndef") {
				end := 0
				isElse := false
				for j := i; j < len(tokens); j++ {
					t := tokens[j]
					if t.Value == "#endif" || t.Value == "#else" {
						if t.Value == "#else" {
							isElse = true
						}
						i = j
						end = j
						break
					}
					out = append(out, t)
				}

				if isElse == true {
					for j := i; j < len(tokens); j++ {
						t := tokens[j]
						if t.Value == "#endif" {	
							i = j
							end = j
							break
						}
					}
				}

				if end == 0 {
					error.Error(29, "", origin, &FakeStream, true)
				}

				again = true
			} else {
				end := 0
				for j := i; j < len(tokens); j++ {
					t := tokens[j]
					if t.Value == "#endif" || t.Value == "#else" {
						i = j
						end = j
						break
					}
				}

				if tokens[i].Value == "#else" {
					i++
					end := 0
					for j := i; j < len(tokens); j++ {
						t := tokens[j]
						if t.Value == "#endif" {
							i = j
							end = j
							break
						}
						out = append(out, t)
					}

					if end == 0 {
						error.Error(29, "", origin, &FakeStream, true)
					}

					again = true
				}

				if end == 0 {
					error.Error(29, "", origin, &FakeStream, true)
				}
			}
		case "#error":
			var FakeStream []shared.Token
			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					break
				}
				FakeStream = append(FakeStream, shared.Token{
					Value: t.Value,
					Line: t.Line,
					File: t.File,
				})
			}
			origin := FakeStream[0]
			i++
			error.Error(22, tokens[i].Value, origin, &FakeStream, true)
		case "#warning":
			var FakeStream []shared.Token
			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					break
				}
				FakeStream = append(FakeStream, shared.Token{
					Value: t.Value,
					Line: t.Line,
					File: t.File,
				})
			}
			origin := FakeStream[0]
			i++
			error.Warning(22, tokens[i].Value, origin, &FakeStream, true)	
		default:
			token := tokens[i]
			if tokens[i].Value[0] == '#' {
				var FakeStream []shared.Token
				for j := i; j < len(tokens); j++ {
					t := tokens[j]
					if t.Value == "\n" {
						break
					}
					FakeStream = append(FakeStream, shared.Token{
						Value: t.Value,
						Line: t.Line,
						File: t.File,
					})
				}
				origin := FakeStream[0]
				error.Error(42, "\"" + token.Value + "\"", origin, &FakeStream, true)
			} else {
				Entry, Found := CheckDefined(token.Value)
				if Found == true {
					for _, Token := range Entry.ReplacementList {
						Token.Line = token.Line
						Token.File = token.File
						out = append(out, Token)
					}
				} else {
					out = append(out, token)
				}
			}
		}
	}

	if again == true {
		tokens = out
		out = []shared.Token {}
		again = false
		goto PREPROCESSOR_TOP
	}

	return out
}
