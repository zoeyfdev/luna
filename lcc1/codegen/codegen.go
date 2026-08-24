package codegen

import (
	"lcc1/neoparser"
	"lcc1/shared"
	"strings"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

var Section1 strings.Builder
var Section2 strings.Builder

func Write(s string, indent bool) {
	indent_s := ""
	if indent == true {
		indent_s = "    "
	}
	Section2.WriteString(indent_s + s + "\n")
}

func WritePre(s string) {
	Section1.WriteString(s + "\n")
}

func CodegenAssignment(Assignment neoparser.Assignment) {
	
}

func CodegenDecl(Decl neoparser.Declaration) {
	switch Decl.(type) {
	case neoparser.Function:
		Func := Decl.(neoparser.Function)

		Name := Func.Name
		Write(Name + ":", false)
		
		// Write prologue
		noreturn := false

		for _, Attr := range Func.Attributes {
			switch Attr {
			case "noreturn":
				noreturn = true
			}
		}

		if noreturn == false {
			Write("pop e11", true)
		}
		
		// TODO: arguments

		Write("push fp", true)
		Write("mov r12, _builtin_lcc_basin_" + Func.Name, true)
		Write("sub fp, fp, r12", true)
		Write("push e11", true)

		// Code goes here

		for _, Statement := range Func.Children {
			switch Statement.(type) {
			case neoparser.Assignment:
				CodegenAssignment(Statement.(neoparser.Assignment))
			}
		}

		// epilogue
		if noreturn == false {
			Write("pop e11", true)
		}
		Write("pop fp", true)
		
		if noreturn == false {
			Write("ret", true)
		}
	}
}

func Codegen(TU *neoparser.AST) string {
	// Reset the buffer
	Section1 = strings.Builder {}
	Section2 = strings.Builder {}

	spew.Dump((*TU))	

	switch shared.Bits {
	case 16:
		WritePre(".bits 16")
	case 32:
		WritePre(".bits 32")
	}

	for _, Decl := range (*TU).Declarations {
		switch Decl.(type) {
		case neoparser.Function:
			Func := Decl.(neoparser.Function)
			if Func.TypeInfo.Static == false {
				WritePre(".global " + Func.Name)
			}	
		}
	}

	WritePre("")

	for _, Decl := range (*TU).Declarations {
		switch Decl.(type) {
		case neoparser.Function:
			Func := Decl.(neoparser.Function)
			WritePre("#define _builtin_lcc_basin_" + Func.Name + " " + fmt.Sprintf("%d", Func.BasinSize))
		}
	}	

	for _, Decl := range (*TU).Declarations {
		CodegenDecl(Decl)
	}

	return Section1.String() + "\n" + Section2.String()
}
