package templates

const test_run_main_temp = `
package main

import (
	"fmt"
	"os"
)

`

//1. 提供一个命令行帮助
const test_build_main_temp = `
func Usage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("%s deploy -- for deploy contract\n", os.Args[0])
	{{range.}}fmt.Printf("%s {{.FuncName}} -- for {{.FuncName}} \n", os.Args[0])
    {{end}}
    fmt.Printf("\n")
}


func main() {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(0)
	}
	if os.Args[1] == "deploy" {
		CallDeploy(os.Args[2:])
	}{{range.}} else if os.Args[1] == "{{.FuncName}}" {
		Call{{.FuncName}}(os.Args[2:])
	} {{end}} else {
		fmt.Printf("params unvalid")
	}
}

`
