package lexer

import (
	"las/shared"
	"strings"
	"unicode"
)

func Tokenize (text string, filename string) []shared.Token {
	var tokens []shared.Token
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
		if len(Value) == 0 { return }
		tokens = append(tokens, shared.Token{Value: Value, Line: Line, File: filename})
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

		if strings.ContainsRune("+-*/%&|^~<>=!?;()[]{}@", r) {
			Write(string(r))
			i++
			continue
		}

		var buf []rune
		for i < len(runes) {
			c := runes[i]
			if unicode.IsSpace(c) || strings.ContainsRune("+-*/%&|^~<>=!?;()[]{}@\"", c) {
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
