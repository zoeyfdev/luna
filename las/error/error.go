package error

import (
	"fmt"
	"os"
)

var Upgrade bool
var Errors int
var Warnings int
var File string
var Line int

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
	"putting more than one character to a register may have undesirable results", // 10 (not used anymore)
	"expected number",
	"unknown pragma directive",
	"deprecated instruction",
	"unknown argument:",
	"invalid hex color to '.fhex'",
}

func Error(errno int, args string) {
	label := ""

	if File != "" {
		label = File + ":"
	} else {
		label = "lcc:"
	}

	if Line > 0 {
		label += fmt.Sprintf("%d:", Line)
	}

	fmt.Fprintln(os.Stderr, "\033[1;39m" + label + " \033[1;31merror: \033[1;39m" + errors[errno] + " " + args + "\033[0m")
	Errors++
}

func Warning(errno int, args string) {
	if Upgrade == true {
		Error(errno, args)
		return
	}
	label := ""

	if File != "" {
		label = File 
	} else {
		label = "lcc"
	}

	fmt.Println("\033[1;39m" + label + " \033[1;35mwarning: \033[1;39m" + errors[errno] + " " + args + "\033[0m")
	Warnings++
}
