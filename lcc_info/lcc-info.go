package lcc_info

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	VERSION string = "7.0"
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
