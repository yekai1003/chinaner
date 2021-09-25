package main

import (
	"chinaner/templates"
	"flag"
	"fmt"
	"log"
	"os"
)

func Usage() {
	fmt.Printf("%s build -file contract'name  -- for compiler contracts\n", os.Args[0])
	fmt.Printf("%s init  -- init test dir\n", os.Args[0])
	fmt.Printf("%s generate -file contract's name  -- generate contract's code\n", os.Args[0])
	fmt.Printf("%s showabi  -- showabi \n", os.Args[0])
}

func run() {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(0)
	}
	buildCMD := flag.NewFlagSet("build", flag.ExitOnError)
	buildCMDName := buildCMD.String("file", "", "contract's name")
	generateCMD := flag.NewFlagSet("generate", flag.ExitOnError)
	generateCMDName := generateCMD.String("file", "", "contract's name")

	jsonCMD := flag.NewFlagSet("json", flag.ExitOnError)
	jsonCMDData := jsonCMD.String("array", "", "contract's params")

	switch os.Args[1] {
	case "build":
		err := buildCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse build CMD")
		}

	case "init":
		//初始化
		if InitDir() != nil {
			log.Panic("failed to init dir")
		}
		fmt.Println("init test dir sucess!")
	case "generate":
		err := generateCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse generate CMD")
		}
	case "json":
		err := jsonCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse json CMD")
		}
	case "showabi":
		templates.ShowAbi("contracts/demo.abi")
	default:
		Usage()
		os.Exit(1)

	}

	if buildCMD.Parsed() {
		//do build code
		if *buildCMDName == "" {
			log.Panic("failed to get contract file name")
		}
		CompilerRun(*buildCMDName)
	}

	if generateCMD.Parsed() {
		//do generate code
		if *generateCMDName == "" {
			log.Panic("failed to get contract file name")
		}
		templates.Run(targetPath, *generateCMDName, "call.go")
	}

	if jsonCMD.Parsed() {
		fmt.Println("json:", *jsonCMDData)
	}

}

func main() {
	run()
}
