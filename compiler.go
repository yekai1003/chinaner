package main

import (
	"chinaner/templates"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gookit/color"
)

//solc --bin --abi --pretty-json -o ./hello  hello.sol
//abigen -sol hello.sol -pkg main -type hello -out hello.go

const targetPath = "contracts"

//编译一个智能合约
func Compiler2Golang(basicName, pkgName string) error {

	//xxx.sol - > xxx.go
	goName := basicName + ".go"
	solName := basicName + ".sol"

	abigenPath := os.Getenv("ABIGEN")
	if abigenPath == "" {
		abigenPath = "abigen"
	}

	cmd := exec.Command(abigenPath, "-sol", solName, "-pkg", pkgName, "-type", basicName, "-out", pkgName+"/"+goName)
	//fmt.Println("Compiler2Golang:", cmd.String())
	color.Green.Println(cmd.String())
	return cmd.Run()
}

//构造abi
func BuildAbi(basicName string) error {
	solcPath := os.Getenv("SOLC")
	if solcPath == "" {
		solcPath = "solc"
	}
	solName := basicName + ".sol"

	cmd := exec.Command(solcPath, "--bin", "--abi", "-o", "./"+targetPath, solName)

	//fmt.Println("BuildAbi:", cmd.String())
	color.Green.Println(cmd.String())
	return cmd.Run()
}

//扫描目录，获得全部的文件
func CompilerRun(contract_name string) error {
	basicName := strings.Replace(contract_name, ".sol", "", -1)
	//fmt.Println("basicname is :", basicName)
	if err := BuildAbi(basicName); err != nil {
		log.Panic("failed to build contract to abi:", err)
	}
	color.Yellow.Println("compile abi sucess!")
	if err := Compiler2Golang(basicName, targetPath); err != nil {
		log.Panic("failed to build contract to golang:", err)
	}
	color.Yellow.Println("compile contracts to golang sucess!")
	color.Yellow.Println("done...")
	return nil
}

//初始化测试目录
func InitDir(dirname string, sdkPath string) error {

	if dirname == "" {
		dirname = "testcontracts"
	}
	cmd := exec.Command("mkdir", "-p", dirname+"/accounts")
	color.Green.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Panic("failed to run:", cmd.String())
	}
	err = templates.InitFile(dirname)
	if err != nil {
		log.Panic("failed to Init file:", err)
	}
	color.Yellow.Println("init cert files sucess!")
	if sdkPath == "" {
		color.Warn.Println("Warn:Please copy ca.crt sdk.crt sdk.key to ", dirname)
	} else {
		cmd = exec.Command("cp", sdkPath+"/ca.crt", dirname)
		color.Green.Println(cmd.String())
		err = cmd.Run()
		if err != nil {
			color.Warn.Println("Warning:sdk path err, please copy ca.crt sdk.crt sdk.key to", dirname)
			return err
		}
		cmd = exec.Command("cp", sdkPath+"/sdk.crt", dirname)
		color.Green.Println(cmd.String())
		err = cmd.Run()
		if err != nil {
			color.Warn.Println("Warning:sdk path err, please copy ca.crt sdk.crt sdk.key to", dirname)
			return err
		}
		cmd = exec.Command("cp", sdkPath+"/sdk.key", dirname)
		color.Green.Println(cmd.String())
		err = cmd.Run()
		if err != nil {
			color.Warn.Println("Warning:sdk path err, please copy ca.crt sdk.crt sdk.key to", dirname)
			return err
		}
	}
	return nil
}

func rungomod(modname string) error {

	cmd := exec.Command("go", "mod", "init", modname)
	color.Green.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Panic("failed to run:", cmd.String())
	}

	cmd = exec.Command("go", "build")
	color.Green.Println(cmd.String())
	err = cmd.Run()
	if err != nil {
		log.Panic("Failed to run:", cmd.String(), err)
	}
	color.Yellow.Println("done...")
	return nil
}
