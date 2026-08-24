package codegen

import (
	"lcc1/neoparser"
	"lcc1/shared"
	"lcc1/error"
	"strings"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
)

var Section1 strings.Builder
var Section2 strings.Builder

var IDCounter int = 1

func Write(s string, indent bool) {
	indent_s :=""
	if indent == true {
		indent_s = "    "
	}
	Section2.WriteString(indent_s + s + "\n")
}

func WritePre(s string) {
	Section1.WriteString(s + "\n")
}

func CodegenLeaf(Leaf neoparser.Leaf) CodegenResult {
	switch Leaf.(type) {
	case neoparser.IntLit:
		IntLit := Leaf.(neoparser.IntLit)
		r := TakeRegister()
		Write("mov " + r + ", " + IntLit.Value, true)
		return CodegenResult {
			Register: r,
			TypeInfo: neoparser.CompositeType {
				Type: neoparser.I8,
			},
		}
	case neoparser.Identifier:
		Identifier := Leaf.(neoparser.Identifier)
		r := TakeRegister()
		Write("mov " + r + ", " + Identifier.AttachedVariable.Internal, true)

		if Identifier.IsRead == true {
			switch Identifier.AttachedVariable.TypeInfo.Type {
			// TODO: add pointers
			case neoparser.I8:
				Write("lod " + r + ", " + r, true)
			case neoparser.I16:
				Write("lod16 " + r + ", " + r, true)
			case neoparser.I32:
				Write("lod32 " + r + ", " + r, true)
			}
		}	

		return CodegenResult {
			Register: r,
			TypeInfo: Identifier.AttachedVariable.TypeInfo,
		}
	}

	return CodegenResult{}
}

func CodegenUnaryOp(UnaryOp neoparser.UnaryOperation) CodegenResult {
	Result := CodegenExpression(UnaryOp.Left)
	switch UnaryOp.Op {
	case shared.TokStar:
		// Dereference
			
		if Result.TypeInfo.PointerLength > 0 {
			Write("lod_ptr " + Result.Register + ", " + Result.Register, true)
			Result.TypeInfo.PointerLength--
		} else {
			switch Result.TypeInfo.Type {
			case neoparser.I8, neoparser.VOID:
				Write("lod " + Result.Register + ", " + Result.Register, true)
			case neoparser.I16:
				Write("lod16 " + Result.Register + ", " + Result.Register, true)
			case neoparser.I32:
				Write("lod32 " + Result.Register + ", " + Result.Register, true)
			}
		}

		return CodegenResult {
			Register: Result.Register,
		}
	default:
		error.InternalCompilerError("invalid unary op type!")
	}

	return CodegenResult {}
}

func CodegenBinaryOp(BinaryOp neoparser.BinaryOperation) CodegenResult {
	Left := CodegenExpression(BinaryOp.Left)
	Right := CodegenExpression(BinaryOp.Right)

	switch BinaryOp.Op {
	case shared.TokPlus:
		Write("add " + Left.Register + ", " + Left.Register + ", " + Right.Register, true)
		FreeRegister(Right.Register)

		return CodegenResult {
			Register: Left.Register,
		}
	case shared.TokMinus:
		Write("sub " + Left.Register + ", " + Left.Register + ", " + Right.Register, true)
		FreeRegister(Right.Register)

		return CodegenResult {
			Register: Left.Register,
		}
	case shared.TokStar:
		Write("mul " + Left.Register + ", " + Left.Register + ", " + Right.Register, true)
		FreeRegister(Right.Register)

		return CodegenResult {
			Register: Left.Register,
		}
	case shared.TokSlash:
		Write("div " + Left.Register + ", " + Left.Register + ", " + Right.Register, true)
		FreeRegister(Right.Register)

		return CodegenResult {
			Register: Left.Register,
		}
	default:
		error.InternalCompilerError("invalid binary op type")
	}


	return CodegenResult {}
}

func CodegenExpression(Expression neoparser.Expression) CodegenResult {
	switch Expression.(type) {
	case neoparser.BinaryOperation:
		return CodegenBinaryOp(Expression.(neoparser.BinaryOperation))
	case neoparser.UnaryOperation:
		return CodegenUnaryOp(Expression.(neoparser.UnaryOperation))
	case neoparser.IntLit, neoparser.Identifier:
		return CodegenLeaf(Expression.(neoparser.Leaf))
	}

	return CodegenResult {}
}

func CodegenStatement(Statement neoparser.Statement) {
	switch Statement.(type) {
	case neoparser.Assignment:
		Assignment := Statement.(neoparser.Assignment)
		Target := CodegenExpression(Assignment.Target)
		Value := CodegenExpression(Assignment.Value)
		
		Write("str " + Target.Register + ", " + Value.Register, true)

		FreeRegister(Target.Register)
		FreeRegister(Value.Register)
	case neoparser.Return:
		Return := Statement.(neoparser.Return)

		Value := CodegenExpression(Return.Value)

		Write("mov e6, " + Value.Register, true)
		Write("ret", true)

		FreeRegister(Value.Register)
	}
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
			CodegenStatement(Statement)	
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

	// spew.Dump((*TU))	

	switch shared.Bits {
	case 16:
		WritePre(".bits 16")
	case 32:
		WritePre(".bits 32")
	}

	for i := 0; i < len((*TU).Declarations); i++ {
		switch (*TU).Declarations[i].(type) {
		case neoparser.Function:
			Func := (*TU).Declarations[i].(neoparser.Function)
			if (*TU).Declarations[i].(neoparser.Function).TypeInfo.Static == false {
				WritePre(".global " + Func.Name)
			}
		case neoparser.Variable:
			Var := (*TU).Declarations[i].(neoparser.Variable)
			Write(Var.Internal + ":", false)
			if Var.Scope == 0 || Var.TypeInfo.Static == true {
				Write(".word 1\n", true)
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
