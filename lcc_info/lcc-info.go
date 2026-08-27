package lcc_info

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	VERSION string = "9.0 (preview)"
	ICE_MESSAGE string = "Please send a bug report to alex@alexflax.xyz or make an issue on the GitHub repository and provide the source code file(s) you used."
)


func PrintVersionInfo() {	
	path, _ := os.Executable()
	path, _ = filepath.EvalSymlinks(path)
	path = filepath.Dir(path)

	fmt.Println("Luna Compiler Collection version " + VERSION)
	fmt.Println("Target: luna-l2")
	fmt.Println("InstalledDir:", path)
}

func PrintUnifiedHelpMessage() {
	fmt.Println(
`OVERVIEW: lcc Luna L2 compiler
USAGE: lcc [options] <file(s)>

OPTIONS:
-help, --help - prints this message
-v - prints the version of LCC currently installed
-o <file> - specifies the output file
-S - specifies to only compile high-level languages to assembly
-c - specifies to only compile and assemble
-define <name> <value> - defines a macro with a given value from the command line
-si - print and execute all commands LCC will use to compile
-scs - print performance statistics about compilation
-sast - specifies to LCC1 to dump the translation unit
-sra - specifies to LCC1 to log all register allocations
-Werror - upgrade all compiler/assembler warnings to errors
-fpie - specifies to all that your program will be compiled in PIE mode
-fpie-16 - specifies to L2LD that it should link your PIE executable in 16-bit mode
-fpie-32 - specifies to L2LD that it should link your PIE executable in 32-bit mode
-fno-autolink - specifies to L2LD to not automatically link to any standard libraries`,
)
	os.Exit(0)
}
