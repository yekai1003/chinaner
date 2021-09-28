package main

import (
	"chinaner/templates"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gookit/color"
)

const version = "1.0.0"
const author = "Yekai"

func Usage() {
	fmt.Printf("\nVersion:\n\t%s\nAuthor:\n\t%s\nUsage:\n", version, author)
	fmt.Printf("\t%s build -file contract's name  -- for compiler contracts\n", os.Args[0])
	fmt.Printf("\t%s init -dir DIRNAME -sdkpath SDKPATH  -- init test dir\n", os.Args[0])
	fmt.Printf("\t%s generate -file contract's name  -- generate contract's code\n", os.Args[0])
	fmt.Printf("\t%s showabi -abi ABIFILE  -- showabi \n\n\n", os.Args[0])
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

	showAbiCMD := flag.NewFlagSet("showabi", flag.ExitOnError)
	showAbiCMDName := showAbiCMD.String("abi", "", "ABIFILE")

	initCMD := flag.NewFlagSet("init", flag.ExitOnError)
	initCMDDir := initCMD.String("dir", "", "DIRNAME")
	initCMDSdkpath := initCMD.String("sdkpath", "", "SDKPATH")

	switch os.Args[1] {
	case "build":
		err := buildCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse build CMD")
		}

	case "init":
		//初始化
		err := initCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse init CMD")
		}

	case "generate":
		err := generateCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse generate CMD")
		}
	case "showabi":

		err := showAbiCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic("faile to Parse showAbi CMD")
		}

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
		rungomod(templates.GetWorkdir())
	}

	if showAbiCMD.Parsed() {

		templates.ShowAbi(*showAbiCMDName)
	}

	if initCMD.Parsed() {
		if InitDir(*initCMDDir, *initCMDSdkpath) != nil {
			log.Panic("failed to init dir")
		}
		color.Yellow.Println("done...")
	}

}

func main() {
	run()
}
