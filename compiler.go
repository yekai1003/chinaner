package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
	fmt.Println("Compiler2Golang:", cmd.String())
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

	fmt.Println("BuildAbi:", cmd.String())
	return cmd.Run()
}

//扫描目录，获得全部的文件
func CompilerRun(contract_name string) error {
	basicName := strings.Replace(contract_name, ".sol", "", -1)
	fmt.Println("basicname is :", basicName)
	if err := BuildAbi(basicName); err != nil {
		log.Panic("failed to build contract to abi:", err)
	}

	if err := Compiler2Golang(basicName, targetPath); err != nil {
		log.Panic("failed to build contract to golang:", err)
	}
	return nil
}

//初始化测试目录
func InitDir() error {
	cmd := exec.Command("git", "clone", "https://gitee.com/teacher233/testcontracts")
	return cmd.Run()
}
