package templates

const test_main_temp = `
package main

import (
	"fmt"
	"log"
	"os"
	"%s/contracts"
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
		fmt.Println("the env:CONTRACT_ADDR  doesn't set")
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
	//fmt.Println({{.OutParams}})
	{{.OutShowRetValue}}

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

const test_pem_info = `-----BEGIN PRIVATE KEY-----
MIGEAgEAMBAGByqGSM49AgEGBSuBBAAKBG0wawIBAQQggLsYjDW2lwBA/ZsSfKGI
v7CXqpE+D7UjEmqq+3dBnluhRANCAATGB/al5rRz38fcRzyUxjFIVOuQuNg+hyaG
UiG+JbXLgBPq0j7TV3hBavEyfYFAf7vivUP3XP7+z+UVqANyqVsS
-----END PRIVATE KEY-----`

const test_pem_public = `-----BEGIN PUBLIC KEY-----
MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAExgf2pea0c9/H3Ec8lMYxSFTrkLjYPocm
hlIhviW1y4AT6tI+01d4QWrxMn2BQH+74r1D91z+/s/lFagDcqlbEg==
-----END PUBLIC KEY-----`

const test_config_info = `[Network]
#type rpc or channel
Type="rpc"
CAFile="ca.crt"
Cert="sdk.crt"
Key="sdk.key"
[[Network.Connection]]
NodeURL="127.0.0.1:8545"
GroupID=1
# [[Network.Connection]]
# NodeURL="127.0.0.1:20200"
# GroupID=2

[Account]
# only support PEM format for now
KeyFile="accounts/0x5946a2ec703e74ce91ac0703396be65daeb5ea99.pem"

[Chain]
ChainID=1
SMCrypto=false`
