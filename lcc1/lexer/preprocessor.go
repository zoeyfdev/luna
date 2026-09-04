package lexer

import (
	"unicode"
	"strings"
	"lcc1/shared"
	"path/filepath"
	"runtime"
	"os"
	"lcc1/error"
)

type SmallToken struct {
	Value string
	Line int
	Filename string
}

func Tokenize (text string, filename string) []SmallToken {
	var tokens []SmallToken
	runes := []rune(text)
	i := 0
	Line := 1

	peek := func(offset int) rune {
		j := i + offset
		if j < len(runes) {
			return runes[j]
		}
		return 0
	}

	Write := func(Value string) {
		tokens = append(tokens, SmallToken{Value: Value, Line: Line, Filename: filename})
	}

	for i < len(runes) {
		r := runes[i]

		if r == '/' && peek(1) == '/' {
			i += 2
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}

		if r == '/' && peek(1) == '*' {
			i += 2
			for i < len(runes) {
				if runes[i] == '*' && peek(1) == '/' {
					i += 2
					break
				}
				if runes[i] == '\n' {
					Line++
				}
				i++
			}
			continue
		}

		if r == '\n' {
			Write("\n")
			Line++
			i++
			continue
		}

		if unicode.IsSpace(r) {
			i++
			continue
		}

		if r == '"' || (r == '<' && tokens[len(tokens) - 1].Value == "#include")  {
			var termonrangle bool
			var buf []rune

			if r == '<' && tokens[len(tokens) - 1].Value == "#include" {
				termonrangle = true
			}

			buf = append(buf, r)
			i++
			for i < len(runes) {
				c := runes[i]
				buf = append(buf, c)
				if c == '\\' && peek(1) != 0 {
					i++
					buf = append(buf, runes[i])
				} else if c == '"' || (c == '>' && termonrangle == true) {
					break
				} else if c == '\n' {
					Line++
				}
				i++
			}
			i++

			Write(string(buf))
			continue
		}

		if r == '=' && peek(1) == '=' {
			Write("==")
			i += 2
			continue
		}

		if r == '>' && peek(1) == '=' {
			Write(">=")
			i += 2
			continue
		}

		if r == '<' && peek(1) == '=' {
			Write("<=")
			i += 2
			continue
		}

		if r == '!' && peek(1) == '=' {
			Write("!=")
			i += 2
			continue
		}

		if r == '+' && peek(1) == '+' {
			Write("++")
			i += 2
			continue
		}

		if r == '-' && peek(1) == '-' {
			Write("--")
			i += 2
			continue
		}

		if r == '>' && peek(1) == '>' {
			Write(">>")
			i += 2
			continue
		}

		if r == '<' && peek(1) == '<' {
			Write("<<")
			i += 2
			continue
		}

		if r == '-' && peek(1) == '>' {
			Write("->")
			i += 2
			continue
		}

		if r == '.' && peek(1) == '.' && peek(2) == '.' {
			Write("...")
			i += 3
			continue
		}

		if strings.ContainsRune("+-*/%&|^~<>=!?:;.,()[]{}@", r) {
			Write(string(r))
			i++
			continue
		}

		var buf []rune
		for i < len(runes) {
			c := runes[i]
			if unicode.IsSpace(c) || strings.ContainsRune("+-*/%&|^~<>=!?:;.,()[]{}@\"", c) {
				break
			}
			buf = append(buf, c)
			i++
		}
		if len(buf) > 0 {
			Write(string(buf))
		}
	}

	return tokens
}

type DefineEntry struct {
	Name string
	ReplacementList []SmallToken
}

var Defines = []DefineEntry {
	DefineEntry { 
		Name: "__LCC__", 
		ReplacementList: []SmallToken { 
			SmallToken { 
				Value: "1", 
				Filename: "__DEFAULT__",
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

func Preprocessor(text string, filename string, predefs []DefineEntry) []SmallToken {
	var out []SmallToken

	Defines = append(Defines, predefs...)

	again := false

	var pragma_once_files []string
	var included_files []string

	CheckPragmaOnce := func(filename string) bool {
		for _, file := range pragma_once_files {
			if file == filename {
				return true
			}
		}
		return false
	}

	tokens := Tokenize(text, filename)
PREPROCESSOR_TOP:
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Value {
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
					Type: shared.TokIdent,
					Value: t.Value,
					FakeValue: t.Value,
					Line: t.Line,
					File: t.Filename,
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
					error.Error(29, "", origin, &FakeStream)
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
						error.Error(29, "", origin, &FakeStream)
					}

					again = true
				}

				if end == 0 {
					error.Error(29, "", origin, &FakeStream)
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
					Type: shared.TokIdent,
					Value: t.Value,
					FakeValue: t.Value,
					Line: t.Line,
					File: t.Filename,
				})
			}
			origin := FakeStream[0]
			i++
			error.Error(22, tokens[i].Value, origin, &FakeStream)
		case "#warning":
			var FakeStream []shared.Token
			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					break
				}
				FakeStream = append(FakeStream, shared.Token{
					Type: shared.TokIdent,
					Value: t.Value,
					FakeValue: t.Value,
					Line: t.Line,
					File: t.Filename,
				})
			}
			origin := FakeStream[0]
			i++
			error.Warning(22, tokens[i].Value, origin, &FakeStream)	
		case "#include":
			again = true
			i++

			raw := ""
			times := 0
			root_by_default := false
			base := filepath.Dir(filename)
			if tokens[i].Value[0] == '"' {
				raw = strings.ReplaceAll(tokens[i].Value, "\"", "")
			} else if tokens[i].Value[0] == '<' {
				raw = strings.ReplaceAll(tokens[i].Value, "<", "")
				raw = strings.ReplaceAll(raw, ">", "")
				root_by_default = true
				times = 1
			}

			path := raw
			if !filepath.IsAbs(raw) {
				path = filepath.Join(base, raw)
			}
			ogpath := path

			if root_by_default == true {
				switch runtime.GOOS {
				case "windows":
					path = "C:\\Program Files (x86)\\Luna L2\\lib\\lcc\\" + raw
				default:
					path = "/usr/local/lib/lcc/" + raw
				}
				ogpath = path
			}

			ff_top:
			if CheckPragmaOnce(path) == true {
				continue
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				times++
				if times >= 2 {
					var FakeStream []shared.Token
					for j := i - 1; j < len(tokens); j++ {
						t := tokens[j]
						if t.Value == "\n" {
							break
						}
						FakeStream = append(FakeStream, shared.Token{
							Type: shared.TokIdent,
							Value: t.Value,
							FakeValue: t.Value,
							Line: t.Line,
							File: t.Filename,
						})
					}
					origin := FakeStream[1]
					error.Error(16, "\"" + ogpath + "\"", origin, &FakeStream)
				} else {
					switch runtime.GOOS {
					case "windows":
						path = "C:\\Program Files (x86)\\Luna L2\\lib\\lcc\\" + raw
					default:
						path = "/usr/local/lib/lcc/" + raw
					} 
					goto ff_top
				}
			}
			i++	

			ntokens := Tokenize(string(contents), path) 
			for j := 0; j < len(ntokens); j++ {
				out = append(out, ntokens[j])
			}

			included_files = append(included_files, path)
		case "#pragma":
			i++
			switch tokens[i].Value {
			case "bits":
				i++
				switch tokens[i].Value {
				case "16":
					shared.Bits = 16
				case "32":
					shared.Bits = 32
				default:
					var FakeStream []shared.Token
					for j := i - 2; j < len(tokens); j++ {
						t := tokens[j]
						if t.Value == "\n" {
							break
						}
						FakeStream = append(FakeStream, shared.Token{
							Type: shared.TokIdent,
							Value: t.Value,
							FakeValue: t.Value,
							Line: t.Line,
							File: t.Filename,
						})
					}
					origin := FakeStream[2]
					error.Warning(43, "", origin, &FakeStream)
				}
			case "once":
				file := tokens[i - 1].Filename
				pragma_once_files = append(pragma_once_files, file)	
			}

			for j := i; j < len(tokens); j++ {
				if tokens[j].Value == "\n" {
					break
				}
				i++
			}
		case "#testcrash":
			error.InternalCompilerError(tokens[i + 1].Value)
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
						Type: shared.TokIdent,
						Value: t.Value,
						FakeValue: t.Value,
						Line: t.Line,
						File: t.Filename,
					})
				}
				origin := FakeStream[0]
				error.Error(42, "\"" + token.Value + "\"", origin, &FakeStream)
			} else {
				Entry, Found := CheckDefined(token.Value)
				if Found == true {
					for _, Token := range Entry.ReplacementList {
						Token.Line = token.Line
						Token.Filename = token.Filename
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
		out = []SmallToken {}
		again = false
		goto PREPROCESSOR_TOP
	}

	return out 
}
