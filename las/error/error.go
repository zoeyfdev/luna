package error

import (
	"fmt"
	"os"
	"las/shared"
	"strings"
)

var Upgrade bool
var Errors int
var Warnings int

var errors = []string {
	"no input files",
	"no such file or directory",
	"invalid register name",
	"invalid operand to instruction",
	"invalid instruction mnemonic",
	"immediate value too large",
	"missing terminating '\"' character",
	"expected string",
	"invalid architecture",
	"invalid argument to 'bits', must be 16 or 32",
	"", // 10 (not used anymore)
	"expected number",
	"unknown pragma directive",
	"deprecated instruction",
	"unknown argument:",
	"invalid hex color to '.fhex'",
	"invalid preprocessing directive",
	"unterminated conditional directive",
}

func StargazeLite(Tokens *[]shared.Token, where int, kind int) {
	// Kinds:
	// 1: error
	// 2: warning
	// 3: note

	line := (*Tokens)[where].Line
	file := (*Tokens)[where].File

	OGTVAL := (*Tokens)[where].Value

	if kind != 3 {
		(*Tokens)[where].Value = "\033[1;31m" + (*Tokens)[where].Value + "\033[0m"
	}

	start := where
	for start > 0 && (*Tokens)[start - 1].Line == line && (*Tokens)[start - 1].File == file {
		start--
	}
	end := where
	for end < len(*Tokens) - 1 && (*Tokens)[end + 1].Line == line && (*Tokens)[end + 1].File == file {
		end++
	}

	words := make([]string, 0, end - start + 1)
	for j := start; j <= end; j++ {
		if (*Tokens)[j].Value == "\n" { continue }
		words = append(words, (*Tokens)[j].Value)	
	}

	text := strings.Join(words, " ")
	fmt.Printf("    %d | %s\n", line, text) 

	if kind != 3 {
		(*Tokens)[where].Value = OGTVAL
	}
}

func find(token shared.Token, tokens *[]shared.Token) int {
	for i, t := range (*tokens) {
		if t == token {
			return i
		}
	}
	return 0
}

func Error(errno int, args string, Token shared.Token, Stream *[]shared.Token, Gaze bool) {
	label := ""

	if Token.File != "" {
		label = Token.File + ":"
	} else {
		label = "lcc:"
	}

	if Token.Line > 0 {
		label += fmt.Sprintf("%d:", Token.Line)
	}

	addtl := " "

	if errno == 10 {
		addtl = ""
	}

	fmt.Fprintln(os.Stderr, "\033[1;39m" + label + " \033[1;31merror: \033[1;39m" + errors[errno] + addtl + args + "\033[0m")

	if Gaze == true {
		StargazeLite(Stream, find(Token, Stream), 1)
	}
	
	Errors++
}

func Warning(errno int, args string, Token shared.Token, Stream *[]shared.Token, Gaze bool) {
	if Upgrade == true {
		Error(errno, args, Token, Stream, Gaze)
		return
	}
	label := ""

	if Token.File != "" {
		label = Token.File 
	} else {
		label = "lcc"
	}

	if Token.Line > 0 {
		label += fmt.Sprintf("%d:", Token.Line)
	}

	addtl := " "

	if errno == 10 {
		addtl = ""
	}

	fmt.Println("\033[1;39m" + label + " \033[1;35mwarning: \033[1;39m" + errors[errno] + addtl + args + "\033[0m")

	if Gaze == true {
		StargazeLite(Stream, find(Token, Stream), 2)
	} else {
		fmt.Printf("\n")
	}
	Warnings++
}
