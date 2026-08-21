package main

import (
	"fmt"
	"os"
	"strings"
	"runtime"
	"os/exec"
	"path/filepath"
	"github.com/alexfdev0/lcc_info"
	"las/error"
	"las/assembler"
	"las/lexer"
	"las/shared"
)

// This code was hell to refactor
// Also it is a bucket of bad so be warned LOL

// this code was aids to refactor AGAIN
// ^ zoey aug 20 2025 日本語授業で

var section string = "text"
var input_files []string
var Current_Filename string = ""

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
		name = name[:len(name)-len(ext)]
	}
	return
}

func cleanupFiles(files []string) {
	if runtime.GOOS != "windows" {
		for _, file := range files {
			execute("rm -f " + file)
		}
	} else {
		for _, file := range files {
			execute("del /f " + file)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		error.Error(0, "", shared.Token{}, &[]shared.Token{}, false)
		os.Exit(1)
	}

	var output_filename string = ""
	var nolink bool = false
	var object_files = []string {}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-v":
			lcc_info.PrintVersionInfo()	
			os.Exit(0)
		case "-o":
			output_filename = os.Args[i + 1]
			i++
		case "-c":
			nolink = true
		case "-Werror":
			error.Upgrade = true
		case "-fpie":
			assembler.PIE = true
		case "-help", "--help":
			lcc_info.PrintUnifiedHelpMessage()
		case "-define":
			if i + 2 > len(os.Args) - 1 {
				continue
			}

			name := os.Args[i + 1]
			value := os.Args[i + 2]

			de := lexer.DefineEntry {}
			
			list := []shared.Token {}

			cur := ""

			for _, c := range value {
				switch c {
				case ' ':
					list = append(list, shared.Token {
						Value: cur,
						File: "__CMDLINE__",
					})
				default:
					cur += string(c)
				}
			}

			de.Name = name
			lexer.Predefs = append(lexer.Predefs, de)

			i += 2
		default:
			if arg[0] == '-' {
				error.Error(14, "'" + arg + "'", shared.Token{}, &[]shared.Token{}, false)
			} else {
				input_files = append(input_files, arg)
			}
		}
	}

	if len(input_files) < 1 {
		error.Error(0, "", shared.Token{}, &[]shared.Token{}, false)
		os.Exit(1)
	}

	if output_filename == "" {
		if nolink == false {
			output_filename = "a.bin"
		} else {
			output_filename = "a.o"
		}
	}

	var link_nocont bool = false

	for _, file := range input_files {
		data, err := os.ReadFile(file)
		if err != nil {
			error.Error(1, "'" + file + "'", shared.Token{}, &[]shared.Token{}, false)
			os.Exit(1)
		}
		Current_Filename = file
		shared.File = file
		// Assemble everything

		
		tokens := lexer.Preprocessor(string(data))
		assembler.Assemble(tokens)

		// Error checking

		var error_str string = ""
		if error.Warnings > 0 {
			error_str = error_str + fmt.Sprintf("%d", error.Warnings) + " warning"
			if error.Warnings > 1 {
				error_str = error_str + "s"
			}
			if error.Errors > 0 {
				error_str = error_str + " and "
			} else {
				error_str = error_str + " generated."
			}
		}
		if error.Errors > 0 {
			error_str = error_str + fmt.Sprintf("%d", error.Errors) + " error"
			if error.Errors > 1 {
				error_str = error_str + "s"
			}
			error_str = error_str + " generated."
		}
		if error.Errors > 0 || error.Warnings > 0 {
			fmt.Println(error_str)
		}
		if error.Errors > 0 {
			link_nocont = true
			continue
		}
		// Write everything
		name, _ := splitFile(file)
		buffer := append(assembler.DataBuffer, assembler.TextBuffer...)
		buffer = append(buffer, assembler.ExtendedDataBuffer...)
		os.WriteFile(name + ".o", buffer, 0644)
		object_files = append(object_files, name + ".o")

		// Reset
		error.Errors = 0
		error.Warnings = 0
		assembler.DataBuffer = []byte {}
		assembler.TextBuffer = []byte {}
		assembler.ExtendedDataBuffer = []byte {}
		section = "text"
	}

	if link_nocont == true {
		os.Exit(1)
	}

	if nolink == true {
		os.Exit(0)
	}

	execute("lcc " + strings.Join(object_files, " ") + " -o " + output_filename)	
	cleanupFiles(object_files)
}
