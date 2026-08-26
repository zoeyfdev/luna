package main

import (
	"fmt"
	"lcc1/error"
	"lcc1/lexer"
	"lcc1/shared"
	"lcc1/parser"
	"lcc1/neoparser"
	"lcc1/codegen"
	"lcc1/typecheck"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alexfdev0/lcc_info"
)

func execute(command string) bool {
	shell := "sh"
	flag := "-c"
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	}

	cmd := exec.Command(shell, flag, command)
	output, err := cmd.CombinedOutput()
	fmt.Printf(string(output))

	if err != nil {
		return false	
	}
	return true
}

func splitFile(path string) (name string, ext string) {
	ext = filepath.Ext(path)
	name = filepath.Base(path)	
	if ext != "" {
		name = name[:len(name) - len(ext)]
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		error.ErrorNoGaze(0, "", shared.Token{Line: 0})
		os.Exit(1)
	}

	var input_files = []string {}
	var output_file string = ""
	var noassemble bool = false
	var nolink bool = false
	var stats bool = false
	var predefs []lexer.DefineEntry
	var OLD bool
	
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-o":
			output_file = os.Args[i + 1]
			i++
		case "-v":
			lcc_info.PrintVersionInfo()	
			os.Exit(0)
		case "-S":
			noassemble = true
		case "-c":
			nolink = true
		case "-Werror":
			error.Upgrade = true
		case "-fpie":
			parser.PIE = true
		case "-scs":
			stats = true
		case "-define":
			if i + 2 > len(os.Args) - 1 {
				continue
			}

			name := os.Args[i + 1]
			value := os.Args[i + 2]

			de := lexer.DefineEntry {}
			
			list := []lexer.SmallToken {}

			cur := ""

			for _, c := range value {
				switch c {
				case ' ':
					list = append(list, lexer.SmallToken {
						Value: cur,
						Filename: "__CMDLINE__",
					})
				default:
					cur += string(c)
				}
			}

			de.Name = name
			predefs = append(predefs, de)

			i += 2
		case "-help", "--help":
			lcc_info.PrintUnifiedHelpMessage()
		case "-old":
			OLD = true
			shared.OLD = true
		default:
			if arg[0] == '-' {
				error.ErrorNoGaze(53, "'" + arg + "'", shared.Token{Line: 0}) 
			} else {
				input_files = append(input_files, arg)
			}
		}
	}

	if output_file == "" {
		if nolink == true {
			output_file = "a.o"
		} else {
			output_file = "a.bin"
		}
	}

	assembly_files := []string {}
	for _, file := range input_files {
		data, err := os.ReadFile(file)
		if err != nil {
			error.ErrorNoGaze(16, "'" + file + "'", shared.Token{Line: 0})
			os.Exit(1)
		}

		if stats == true {
			fmt.Println("\033[1;39mlcc: statistics in file '" + file + "':\033[0m")
		}

		preproc_time := time.Now()
		code := lexer.Preprocessor(string(data), file, predefs)

		if stats == true {
			fmt.Println("\033[1;39mlcc: preprocessor:", time.Since(preproc_time), "\033[0m")
		}

		lex_time := time.Now()
		tokens := lexer.Lex(code, file)

		if stats == true {
			fmt.Println("\033[1;39mlcc: lexer:", time.Since(lex_time), "\033[0m")
		}

		parse_time := time.Now()
		Code := ""


		if OLD == false {
			AST := neoparser.Parse(tokens)
			typecheck.TypeCheck(AST)

			if error.Errors <= 0 {
				Code = codegen.Codegen(AST)
			}
		} else {
			parser.Parse(tokens, 1)
		}

		if stats == true {
			fmt.Println("\033[1;39mlcc: parser:", time.Since(parse_time), "\033[0m")
		}

		if error.Errors < 1 && error.FailCompilation == false {
			name := strings.TrimSuffix(file, filepath.Ext(file))
			assembly_files = append(assembly_files, name + ".s")

			if OLD == false {
				os.WriteFile(name + ".s", []byte(Code), 0644)
			} else {
				dir := ""

				switch shared.Bits {
				case 16:
					dir = ".bits 16"
				case 32:
					dir = ".bits 32"
				default:
					dir = ".bits 16"
				}

				for _, variable := range parser.Variables {
					if variable.HasBasin == true {
						parser.Code1 = "#define _builtin_lcc_basin_" + variable.Name + " " + fmt.Sprintf("%d", variable.BasinSize) + "\n" + parser.Code1
					}
				}

				os.WriteFile(name + ".s", []byte(dir + "\n" + parser.Code1 + "\n" + parser.Code2), 0644)
			}
		}	
	}

	error.Summary()

	if error.Errors > 0 {
		os.Exit(1)
	}

	if noassemble == true {
		os.Exit(0)
	}

	if nolink == true {
		execute("lcc -c " + strings.Join(assembly_files, " "))
	} else {
		execute("lcc " + strings.Join(assembly_files, " ") + " -o " + output_file)
	}

	for _, file := range assembly_files {
		if runtime.GOOS != "windows" {
			execute("rm -f " + file)
		} else {
			execute("del /f " + file)
		}
	}
}
