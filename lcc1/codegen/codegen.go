package codegen

import (
	"lcc1/neoparser"
	"lcc1/shared"
	"lcc1/error"
	"strings"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

var Section1 strings.Builder
var Section2 strings.Builder

var ReturnLabel string

var IDCounter int = 1

var DumpTU bool = false

func Write(s string, indent bool) {
	indent_s :=""
	if indent == true {
		indent_s = "    "
	}
	Section2.WriteString(indent_s + s + "\n")
}

func WritePre(s string, indent bool) {
	indent_s :=""
	if indent == true {
		indent_s = "    "
	}
	Section1.WriteString(indent_s + s + "\n")
}

func CodegenLeaf(Leaf neoparser.Leaf) CodegenResult {
	switch Leaf.(type) {
	case neoparser.IntLit:
		IntLit := Leaf.(neoparser.IntLit)
		r := TakeRegister()
		Result := CodegenResult {
			Register: r,
			TypeInfo: neoparser.CompositeType {
				Type: ReturnUintPtrType(),
			},
			IsRvalue: true,
		}

		/*
		if IsRead == true {
			Result.Read = true
		}
		*/

		Result.Read = true	
		Write("mov " + r + ", " + IntLit.Value, true)
		return Result
	case neoparser.StringLit:
		StringLit := Leaf.(neoparser.StringLit)
		r := TakeRegister()
		Result := CodegenResult {
			Register: r,
			TypeInfo: StringLit.Type,
			OriginalPointerLength: StringLit.Type.PointerLength,
		} 

		StringName := fmt.Sprintf("var_str_%d", VarTicker)
		WritePre(StringName + ":", false)
		WritePre(".asciz \"" + StringLit.Value + "\"\n", true)
		VarTicker++
		
		LabelName := fmt.Sprintf("var_ptr_%d", VarTicker)
		WritePre(LabelName + ":", false)
		WritePre(".ptr " + StringName + "\n", true)
		VarTicker++
	
		Write("mov " + r + ", " + LabelName, true)

		if StringLit.IsRead == true {
			Result.Read = true
			Write("lod_ptr " + r + ", " + r, true)	
		}

		return Result
	case neoparser.Identifier:
		Identifier := Leaf.(neoparser.Identifier)	
		r := TakeRegister()
		Result := CodegenResult {
			Register: r,
			TypeInfo: Identifier.Type,
			OriginalPointerLength: Identifier.Type.PointerLength,
		}

		Write("mov " + r + ", " + Identifier.AttachedVariable.Internal, true)

		if Identifier.IsRead == true && Identifier.AttachedVariable.Register == false  {
			Result.Read = true
			if Identifier.Type.PointerLength > 0 {
				Write("lod_ptr " + r + ", " + r, true)
				Result.TypeInfo.PointerLength--
			} else {
				switch Identifier.Type.Type {
				// TODO: add pointers
				case neoparser.I8:
					Write("lod " + r + ", " + r, true)
				case neoparser.I16:
					Write("lod16 " + r + ", " + r, true)
				case neoparser.I32:
					Write("lod32 " + r + ", " + r, true)
				}
			}	
		}

		return Result
	}

	return CodegenResult {}
}

func CodegenUnaryOp(UnaryOp neoparser.UnaryOperation, IsWrite bool) CodegenResult {
	// TODO: set read true for this
	switch UnaryOp.Op {
	case shared.TokStar:
		// Dereference	
		Result := CodegenExpression(UnaryOp.Left, IsWrite)

		if IsWrite == true && Result.IsRvalue == true && Result.RValueDerefs <= 0 {
			Write("// No load for you", true)
			Result.RValueDerefs++
			return Result
		}

		Write("// Star!", true)

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
			default:
				error.InternalCompilerError("couldn't dereference type " + fmt.Sprintf("%d", Result.TypeInfo.Type))
			}
		}

		return Result
	case shared.TokAmpersand:
		UnaryOp.Left = neoparser.SetRead(UnaryOp.Left, false)
		Result := CodegenExpression(UnaryOp.Left, IsWrite)

		return Result
	default:
		error.InternalCompilerError("invalid unary op type!")
	}

	return CodegenResult {}
}

func CodegenBinaryOp(BinaryOp neoparser.BinaryOperation, IsWrite bool) CodegenResult {
	BinaryOp.Left = neoparser.SetRead(BinaryOp.Left, true)
	BinaryOp.Right = neoparser.SetRead(BinaryOp.Right, true)

	Left := CodegenExpression(BinaryOp.Left, IsWrite)
	Right := CodegenExpression(BinaryOp.Right, IsWrite)

	OpMap := make(map[shared.TokenType]string)
	OpMap[shared.TokPlus] = "add"
	OpMap[shared.TokMinus] = "sub"
	OpMap[shared.TokStar] = "mul"
	OpMap[shared.TokSlash] = "div"
	OpMap[shared.TokPercent] = "mod"

	OpMap[shared.TokEquality] = "cmp"
	OpMap[shared.TokInequality] = "cmp"
	OpMap[shared.TokLAngle] = "ilt"
	OpMap[shared.TokRAngle] = "igt"
	OpMap[shared.TokGEqual] = "ilt"
	OpMap[shared.TokLEqual] = "igt"

	if OpMap[BinaryOp.Op] == "" {
		println(BinaryOp.Token.Line)
		error.InternalCompilerError("invalid binary operation, got " + fmt.Sprintf("%d", BinaryOp.Op))
	}

	Write(OpMap[BinaryOp.Op] + " " + Left.Register + ", " + Left.Register + ", " + Right.Register, true)

	switch BinaryOp.Op {
	case shared.TokInequality, shared.TokGEqual, shared.TokLEqual:
		r := TakeRegister()
		Write("mov " + r + ", 1", true)
		Write("xor " + Left.Register + ", " + Left.Register + ", " + r, true)
		FreeRegister(r)	
	}

	FreeRegister(Right.Register)
		
	Result := TypeMediation(Left, Right)	

	Result.Register = Left.Register
	Result.IsRvalue = true
	return Result
}

func CodegenCast(Cast neoparser.Cast, IsWrite bool) CodegenResult {
	Value := CodegenExpression(Cast.Value, IsWrite)	
	Value.TypeInfo = Cast.Type

	if Value.Read == true {
		Value.TypeInfo.PointerLength--
	}

	return Value
}

func CodegenExpression(Expression neoparser.Expression, IsWrite bool) CodegenResult {
	switch Expression.(type) {
	case neoparser.BinaryOperation:
		return CodegenBinaryOp(Expression.(neoparser.BinaryOperation), IsWrite)
	case neoparser.UnaryOperation:
		return CodegenUnaryOp(Expression.(neoparser.UnaryOperation), IsWrite)
	case neoparser.IntLit, neoparser.Identifier, neoparser.StringLit:
		return CodegenLeaf(Expression.(neoparser.Leaf))
	case neoparser.Cast:
		return CodegenCast(Expression.(neoparser.Cast), IsWrite)
	case neoparser.IncrementDecrement:
		IncrementDecrement := Expression.(neoparser.IncrementDecrement)

		// TODO: pre-increment/post-increment

		_Value := IncrementDecrement.Target
		IncrementDecrement.Target = neoparser.SetRead(IncrementDecrement.Target, false)
		_Value = neoparser.SetRead(_Value, true)

		r := TakeRegister()
		Target := CodegenExpression(IncrementDecrement.Target, true)
		Value := CodegenExpression(_Value, false)
		
		Write("mov " + r + ", " + Value.Register, true)
		
		if IncrementDecrement.Decrement == false {
			Write("inc " + r, true)
		} else {
			Write("dec " + r, true)
		}

		op := "str"
		if Target.TypeInfo.PointerLength > 0 {
			op = "str_ptr"
		} else {
			switch Target.TypeInfo.Type {
			case neoparser.I8:
				op = "str"
			case neoparser.I16:
				op = "str16"
			case neoparser.I32:
				op = "str32"
			}
		}

		Write(op + " " + Target.Register + ", " + r, true)

		FreeRegister(Target.Register)
		FreeRegister(r)

		return Value
	case neoparser.Assignment:
		Assignment := Expression.(neoparser.Assignment)
		Target := CodegenExpression(Assignment.Target, true)
		Value := CodegenExpression(Assignment.Value, false)

		op := "str"
		if Target.TypeInfo.PointerLength > 0 {
			op = "str_ptr"
		} else {
			switch Target.TypeInfo.Type {
			case neoparser.I8:
				op = "str"
			case neoparser.I16:
				op = "str16"
			case neoparser.I32:
				op = "str32"
			}
		}

		Write(op + " " + Target.Register + ", " + Value.Register, true)

		FreeRegister(Target.Register)

		return Value
	case neoparser.FunctionCall:
		FunctionCall := Expression.(neoparser.FunctionCall)

		Write("// Push allocated registers", true)
		PushAllocated()

		for _, Child := range FunctionCall.Children {
			Value := CodegenExpression(Child, false)
			Write("push " + Value.Register, true)
			FreeRegister(Value.Register)
		}
	
		Write("call " + FunctionCall.AttachedVariable.Internal, true)
		Write("// Pop saved registers", true)
		PopAllocated()

		r := TakeRegister()
		Write("mov " + r + ", " + "e6", true)

		return CodegenResult {
			Register: r,
			IsRvalue: true,
		}
	case neoparser.StructAccess:
		StructAccess := Expression.(neoparser.StructAccess)
		Target := CodegenExpression(StructAccess.Target, IsWrite)

		Write("// Above should not be loaded", true)

		Result := CodegenResult {
			Register: Target.Register,
			TypeInfo: StructAccess.Type,
		}
		
		if StructAccess.Pointer == true && Target.IsRvalue == false {
			Write("lod_ptr " + Target.Register + ", " + Target.Register, true)
		}

		r := TakeRegister()
		Write("mov " + r + ", " + fmt.Sprintf("%d", StructAccess.Offset), true)
		Write("add " + Target.Register + ", " + Target.Register + ", " + r, true)	

		if StructAccess.IsRead == true {
			Result.Read = true
			if StructAccess.Type.PointerLength > 0 {
				Write("lod_ptr " + Target.Register + ", " + Target.Register, true)
				Result.TypeInfo.PointerLength--
			} else {
				switch StructAccess.Type.Type {
				// TODO: add pointers
				case neoparser.I8:
					Write("lod " + Target.Register + ", " + Target.Register, true)
				case neoparser.I16:
					Write("lod16 " + Target.Register + ", " + Target.Register, true)
				case neoparser.I32:
					Write("lod32 " + Target.Register + ", " + Target.Register, true)
				}
			}
		}
		
		FreeRegister(r)
		
		return Result 
	}

	return CodegenResult {}
}

func CodegenStatement(Statement neoparser.Statement) {
	switch Statement.(type) {
	case neoparser.Assignment:
		FreeRegister(CodegenExpression(Statement.(neoparser.Expression), false).Register)
	case neoparser.FunctionCall:
		FreeRegister(CodegenExpression(Statement.(neoparser.Expression), false).Register)
	case neoparser.StatementExpression:
		StatementExpression := Statement.(neoparser.StatementExpression)
		FreeRegister(CodegenExpression(StatementExpression.Expression, false).Register)
	case neoparser.Return:
		Return := Statement.(neoparser.Return)
		Value := CodegenExpression(Return.Value, false)

		Write("mov e6, " + Value.Register, true)
		Write("jmp " + ReturnLabel, true)

		FreeRegister(Value.Register)
	case neoparser.IfStatement:
		IfStatement := Statement.(neoparser.IfStatement)
		VT := VarTicker
		VarTicker++

		Write(fmt.Sprintf("if_stmt_%d_check:", VT), false)

		Result := CodegenExpression(IfStatement.Condition, false)
		Write("jnz " + Result.Register + ", " + fmt.Sprintf("if_stmt_%d_success", VT), true)
		Write("jmp " + fmt.Sprintf("if_stmt_%d_else", VT), true)

		FreeRegister(Result.Register)

		Write(fmt.Sprintf("if_stmt_%d_success:", VT), false)
		for _, Stmt := range IfStatement.SuccessChildren {
			CodegenStatement(Stmt)
		}
		Write("jmp " + fmt.Sprintf("if_stmt_%d_done", VT), true)
		Write(fmt.Sprintf("if_stmt_%d_else:", VT), false)
		if len(IfStatement.ElseChildren) > 0 {
			for _, Stmt := range IfStatement.ElseChildren {
				CodegenStatement(Stmt)
			}
		}
		Write(fmt.Sprintf("if_stmt_%d_done:", VT), false)
	case neoparser.WhileStatement:
		WhileStatement := Statement.(neoparser.WhileStatement)
		VT := VarTicker
		VarTicker++

		CheckLabel := fmt.Sprintf("while_stmt_%d_check", VT)
		BodyLabel := fmt.Sprintf("while_stmt_%d_body", VT)
		AfterLabel := fmt.Sprintf("while_stmt_%d_after", VT)

		OCCL := CurrentContinueLabel
		CurrentContinueLabel = CheckLabel
		OCEL := CurrentEndLabel
		CurrentEndLabel = AfterLabel

		Write(CheckLabel + ":", false)

		Result := CodegenExpression(WhileStatement.Condition, false)
		Write("jnz " + Result.Register + ", " + BodyLabel, true)
		Write("jmp " + AfterLabel, true)
		FreeRegister(Result.Register)

		Write(BodyLabel + ":", false)
		for _, Statement := range WhileStatement.Children {
			CodegenStatement(Statement)
		}
		Write("jmp " + CheckLabel, true)

		Write(AfterLabel + ":", false)

		CurrentEndLabel = OCEL
		CurrentContinueLabel = OCCL
	case neoparser.ForStatement:
		ForStatement := Statement.(neoparser.ForStatement)
		VT := VarTicker
		VarTicker++

		InitLabel := fmt.Sprintf("for_stmt_%d_init", VT) 
		CheckLabel := fmt.Sprintf("for_stmt_%d_check", VT)
		BodyLabel := fmt.Sprintf("for_stmt_%d_body", VT)
		IteratorLabel := fmt.Sprintf("for_stmt_%d_iterator", VT)
		AfterLabel := fmt.Sprintf("for_stmt_%d_after", VT)

		OCCL := CurrentContinueLabel
		CurrentContinueLabel = IteratorLabel
		OCEL := CurrentEndLabel
		CurrentEndLabel = AfterLabel

		Write(InitLabel + ":", false)
		for _, Statement := range ForStatement.Initializer {
			CodegenStatement(Statement)
		}

		Write(CheckLabel + ":", false)
		Result := CodegenExpression(ForStatement.Condition, false)
		Write("jnz " + Result.Register + ", " + BodyLabel, true)
		Write("jmp " + AfterLabel, true)

		Write(BodyLabel + ":", false)
		for _, Statement := range ForStatement.Children {
			CodegenStatement(Statement)
		}

		Write(IteratorLabel + ":", false)
		CodegenExpression(ForStatement.Iterator, false)
		Write("jmp " + CheckLabel, true)

		Write(AfterLabel + ":", false)

		CurrentEndLabel = OCEL
		CurrentContinueLabel = OCCL
	case neoparser.Assembly:
		Assembly := Statement.(neoparser.Assembly)
		Write(Assembly.String + "    // User-defined inline assembly", true)
	case neoparser.GotoStatement:
		GotoStatement := Statement.(neoparser.GotoStatement)
		Write("jmp " + GotoStatement.Name, true)
	case neoparser.BreakStatement:
		Write("jmp " + CurrentEndLabel, true)
	case neoparser.ContinueStatement:
		Write("jmp " + CurrentContinueLabel, true)
	}
}

func CodegenDecl(Decl neoparser.Declaration) {
	switch Decl.(type) {
	case neoparser.Variable:
		Var := Decl.(neoparser.Variable)
		switch Var.Kind {
		case neoparser.FUNCTION:
			if Var.TypeInfo.Extern == false {
				Name := Var.Internal
				ReturnLabel = "." + Name + "_ret"

				Write(Name + ":", false)
				
				// Prologue
				noreturn := false

				for _, Attr := range Var.Attributes {
					switch Attr {
					case "noreturn":
						noreturn = true
					}
				}

				if noreturn == false {
					Write("pop e11", true)
				}	

				// Pop arguments to registers and save them
				r := len(Var.Parameters) - 1
				for i := r; i >= 0; i-- {
					Write(fmt.Sprintf("pop e%d", r), true)
					r--
				}	

				Write("push fp", true)
				Write("mov r12, _builtin_lcc_basin_" + Var.Internal, true)
				Write("sub fp, fp, r12", true)

				if noreturn == false {
					Write("push e11", true)
				}

				r = 0
				for i := r; r < len(Var.Parameters); i++ {
					Param := Var.Parameters[i]
					register := fmt.Sprintf("e%d", r)

					Write("mov r0, " + Param.Internal, true)

					if Param.TypeInfo.PointerLength > 0 || Param.TypeInfo.Type == neoparser.VOID {
						Write("str_ptr r0, " + register, true)
					} else {
						switch Param.TypeInfo.Type {
						case neoparser.I8:
							Write("str r0, " + register, true)
						case neoparser.I16:
							Write("str16 r0, " + register, true)
						case neoparser.I32:
							Write("str32 r0, " + register, true)
						}	
					}

					r++
				}

				for _, Statement := range Var.Children {
					CodegenStatement(Statement)	
				}

				// Epilogue	
				Write("." + Name + "_ret:", false)
				if noreturn == false {
					Write("pop e11", true)
				}
				Write("pop fp", true)
				
				if noreturn == false {
					Write("ret", true)
				}	

				Write("", false)
			}	
		case neoparser.VARIABLE:
			if (Var.Scope == 0 || Var.TypeInfo.Static == true) && Var.Register == false {
				Value := ""
				if Var.TypeInfo.Extern == false {
					Value = Const_CodegenExpression(Var.Children[0].(neoparser.ConstAssignStatement).Value, false)
				} else {
					Value = Var.Name
				}
				WritePre(Var.Internal + ":", false)
				if Var.TypeInfo.PointerLength > 0 {
					WritePre(".ptr " + Value, true)
				} else {
					switch Var.TypeInfo.Type {
					case neoparser.I8:
						WritePre(".byte " + Value, true)
					case neoparser.I16:
						WritePre(".word " + Value, true)
					case neoparser.I32:
						WritePre(".dword " + Value, true)
					}
				}

				WritePre("", false)
			} else {
				// Local variable decls...
			}
		}
	case neoparser.Assembly:
		Assembly := Decl.(neoparser.Assembly)
		Write(Assembly.String + "    // User-defined inline assembly", false)
	}
}

func Codegen(TU *neoparser.AST) string {
	// Reset the buffer
	Section1 = strings.Builder {}
	Section2 = strings.Builder {}

	if DumpTU == true {
		error.NoteCustom("current translation unit:")
		spew.Dump((*TU))
	}

	switch shared.Bits {
	case 16:
		WritePre(".bits 16", false)
	case 32:
		WritePre(".bits 32", false)
	}

	for i := 0; i < len((*TU).Declarations); i++ {
		switch (*TU).Declarations[i].(type) {
		case neoparser.Variable:
			Var := (*TU).Declarations[i].(neoparser.Variable)
			switch Var.Kind {
			case neoparser.FUNCTION:
				if Var.TypeInfo.Static == false && Var.TypeInfo.Extern == false {
					WritePre(".global " + Var.Name, false)
				}	
			}	
		}
	}
	
	WritePre("", false)

	for _, Decl := range (*TU).Declarations {
		switch Decl.(type) {
		case neoparser.Variable:
			Variable := Decl.(neoparser.Variable)
			switch Variable.Kind {
			case neoparser.FUNCTION:
				if Variable.TypeInfo.Extern == false {
					WritePre("#define _builtin_lcc_basin_" + Variable.Internal + " " + fmt.Sprintf("%d", Variable.BasinSize), false)
				}
			}
		}
	}

	WritePre("", false)

	for _, Decl := range (*TU).Declarations {
		CodegenDecl(Decl)
	}

	return Section1.String() + "\n" + Section2.String()
}
