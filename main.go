package main

import (
	//"chinaner/templates"
	"flag"
	"fmt"
	"log"
	"os"
)

func Usage() {
	fmt.Printf("%s build -file contract'name  -- for compiler contracts\n", os.Args[0])
	fmt.Printf("%s 2  -- build test code\n", os.Args[0])
}

func run() {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(0)
	}
	buildCMD := flag.NewFlagSet("build", flag.ExitOnError)
	contractNname := buildCMD.String("file", "", "contract's name")

	if os.Args[1] == "build" {
		err := buildCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse build CMD")
		}
	}

	if os.Args[1] == "deploy" {

	}

	if buildCMD.Parsed() {
		//do build code
		if *contractNname == "" {
			log.Panic("failed to get contract file name")
		}
		fmt.Println("We need to build code", *contractNname, "----")

	}

}

func main() {
	run()
}

// func main() {

// 	if os.Args[1] == "1" {
// 		CompilerRun()
// 	} else if os.Args[1] == "2" {
// 		//build test code
// 		templates.Run()
// 	} else {
// 		Usage()
// 		os.Exit(0)
// 	}

// }
