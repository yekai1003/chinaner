package templates

const test_main_temp = `
package main

import (
	"fmt"
	"log"
	"os"
	"testcontracts/contracts"
	%s
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/ethereum/go-ethereum/common"
)

var testClient *client.Client
var contract_addr string

func init() {
	//1. 解析配置文件
	configs, err := conf.ParseConfigFile("config.toml")
	if err != nil {
		log.Fatalf("ParseConfigFile failed, err: %v", err)
	}
	//2. 连接节点
	client, err := client.Dial(&configs[0])
	if err != nil {
		log.Fatal(err)
	}
	testClient = client

	contract_addr = os.Getenv("CONTRACT_ADDR")
}

`
const test_deploy_temp = `
func CallDeploy(params []string) error {
{{.PrepareParams}}
	addr, ts, _, err := contracts.{{.DeployName}}{{.DeployParams}}
	if err != nil {
		fmt.Println("failed to Deploy contracts", err)
		return err
	}

	fmt.Println("addr=", addr.Hex(), "hash=", ts.Hash().Hex())
	return err
}
`

const test_nogas_temp = `
func Call{{.FuncName}}(params []string) error {
{{.PrepareParams}}
	if contract_addr == "" {
		fmt.Println("the contract_addr doesn't set")
		return nil
	}
	//使用之前部署得到的合约地址
	instance, err := contracts.{{.NewContractName}}(common.HexToAddress(contract_addr), testClient)
	if err != nil {
		fmt.Println("failed to instance contract", err)
		return err
	}
	//调用合约函数
	{{.OutParams}} := instance.{{.FuncName}}{{.InputParams}}
	fmt.Println({{.OutParams}})

	return err
}
`

const test_gas_temp = `
func Call{{.FuncName}}(params []string) error {
	
{{.PrepareParams}}
	if contract_addr == "" {
		fmt.Println("the contract_addr doesn't set")
		return nil
	}
	//1. 构造函数入口 - 合约对象
	instance, err := contracts.{{.NewContractName}}(common.HexToAddress(contract_addr), testClient)
	if err != nil {
		fmt.Println("failed to contract instance", err)
		return err
	}
	//2. 调用函数

	ts, _, err := instance.{{.FuncName}}{{.InputParams}}
	if err != nil {
		fmt.Println("failed to Deposit ", err)
		return err
	}
	fmt.Println(ts.Hash().Hex())
	return err
}
`
