package codegen

import (
	"lcc1/error"
)

type Register struct {
	Name string
	Taken bool
}

var Registers = []Register {
	Register { Name: "r1", Taken: false, },
	Register { Name: "r2", Taken: false, },
	Register { Name: "r3", Taken: false, },
	Register { Name: "r4", Taken: false, },
	Register { Name: "r5", Taken: false, },
	Register { Name: "r6", Taken: false, },
	Register { Name: "r7", Taken: false, },
	Register { Name: "r8", Taken: false, },
	Register { Name: "r9", Taken: false, },
	Register { Name: "r10", Taken: false, },
	Register { Name: "r11", Taken: false, },
	Register { Name: "r12", Taken: false, },
	Register { Name: "e7", Taken: false, },
	Register { Name: "e8", Taken: false, },
	Register { Name: "e9", Taken: false, },
	Register { Name: "e10", Taken: false, },
	Register { Name: "e11", Taken: false, },
	Register { Name: "e12", Taken: false, },
}

var ShowRegisterAllocation bool = false
func TakeRegister() string {
	for i := 0; i < len(Registers); i++ {
		if Registers[i].Taken == false {
			Registers[i].Taken = true
			
			if ShowRegisterAllocation == true {
				error.NoteCustom("allocating register " + Registers[i].Name)
			}

			return Registers[i].Name
		}
	}

	error.InternalCompilerError("no more free registers!")
	return ""
}

func FreeRegister(name string) {
	for i := 0; i < len(Registers); i++ {
		if Registers[i].Name == name {

			if ShowRegisterAllocation == true {
				error.NoteCustom("freeing register " + Registers[i].Name)
			}

			Registers[i].Taken = false
		}
	}
}

func PushAllocated() {
	for _, Register := range Registers {
		if Register.Taken == true {
			Write("push " + Register.Name, true)
		}
	}
}

func PopAllocated() {
	for i := len(Registers) - 1; i >= 0; i-- {
		Register := Registers[i]
		if Register.Taken == true {
			Write("pop " + Register.Name, true)
		}
	}
}
