package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

type DeployContractParams struct {
	DeployName    string
	DeployParams  string
	PrepareParams string
}

//无gas函数调用
type FuncNoGasParams struct {
	FuncName        string
	NewContractName string
	OutParams       string
	InputParams     string
	PrepareParams   string
}

//有gas函数调用
type FuncGasParams struct {
	FuncName        string
	NewContractName string
	InputParams     string
	PrepareParams   string
}

type InputsOutPuts struct {
	Name string
	Type string
}

type FuncInfo struct {
	FuncName string
	Num      int
}

type AbiInfo struct {
	Constant        bool
	Inputs          []InputsOutPuts
	Name            string
	Outputs         []InputsOutPuts
	Payable         bool
	StateMutability string
	Type            string
}

func isHasUint(infos []AbiInfo) bool {
	for _, v := range infos {
		for _, k := range v.Inputs {
			if strings.Contains(k.Type, "uint") {
				return true
			}
		}
	}

	return false
}

func readAbi(abifile string) ([]AbiInfo, error) {
	file, err := os.Open(abifile)
	if err != nil {
		fmt.Println("failed to open file ", err)
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("failed to read abi", err)
		return nil, err
	}
	var abiInfos []AbiInfo
	strdata := strings.Replace(string(data), "\\", "", -1)
	err = json.Unmarshal([]byte(strdata), &abiInfos)
	if err != nil {
		fmt.Println("failed to Unmarshal abi", err)
		return nil, err
	}
	return abiInfos, err
}

func Impl_run_code(runCodePath, runCodeName, basicName string) error {
	//1. 写到哪
	outfile, err := os.OpenFile(runCodeName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile.Close()

	// 读取abi文件信息
	abiInfos, err := readAbi(runCodePath + "/" + basicName + ".abi")
	if err != nil {
		fmt.Println("failed to read abi", err)
		return err
	}
	//2. 写什么
	maintempdata := test_main_temp
	if isHasUint(abiInfos) {
		maintempdata = fmt.Sprintf(test_main_temp, "\t\"strconv\"\n\t\"math/big\"")
	}

	_, err = outfile.WriteString(maintempdata)
	if err != nil {
		fmt.Println("failed to write ", err)
		return err
	}

	//fmt.Println(infos)

	//3. 写入部署合约代码
	//定义部署模版
	deploy_temp, err := template.New("deploy").Parse(test_deploy_temp)
	if err != nil {
		fmt.Println("failed to template deploy", err)
		return err
	}
	var deploy_data DeployContractParams
	deploy_data.DeployName = "Deploy" + strings.Title(basicName)

	//定义nogas函数的模版
	nogas_temp, err := template.New("nogas").Parse(test_nogas_temp)
	if err != nil {
		fmt.Println("failed to template nogas_temp", err)
		return err
	}

	var func_nogas_data FuncNoGasParams
	func_nogas_data.NewContractName = "New" + strings.Title(basicName)

	//定义有gas模版
	hasgas_temp, err := template.New("hasgas").Parse(test_gas_temp)
	if err != nil {
		fmt.Println("failed to template hasgas_temp", err)
		return err
	}

	var func_gas_data FuncGasParams
	func_gas_data.NewContractName = "New" + strings.Title(basicName)

	//对abi进行遍历处理
	for _, v := range abiInfos {
		v.Name = strings.Title(v.Name) //标题优化，首字母大写, hello world - > Hello World
		prepareData := makePrepareParams(v.Inputs)
		deploy_data.PrepareParams = prepareData
		if v.Type == "constructor" {

			// 如果是构造函数-部署函数
			deploy_data.DeployParams = "(testClient.GetTransactOpts(),testClient"
			for num, vv := range v.Inputs {
				//需要根据输入数据类型来判断如何处理:string,address,uint256
				paramName := fmt.Sprintf("param%d", num)
				if vv.Type == "address" {
					deploy_data.DeployParams += " ," + paramName
				} else if vv.Type == "uint256" {
					deploy_data.DeployParams += fmt.Sprintf(" ,%s", paramName)
				} else if vv.Type == "string" {
					deploy_data.DeployParams += " ," + paramName
				} else if vv.Type == "uint" {
					deploy_data.DeployParams += fmt.Sprintf(" ,%s", paramName)
				} else {
					deploy_data.DeployParams += fmt.Sprintf(" ,%s", paramName)
				}

			}
			deploy_data.DeployParams += ")"
			//模版的执行
			err = deploy_temp.Execute(outfile, &deploy_data)
			if err != nil {
				fmt.Println("failed to template Execute ", err)
				return err
			}
		} else {
			//处理其他函数
			if v.StateMutability == "view" {
				//不需要gas函数
				func_nogas_data.FuncName = v.Name
				func_nogas_data.PrepareParams = prepareData
				func_nogas_data.InputParams = "(testClient.GetCallOpts()"

				for num, vv := range v.Inputs {
					//需要根据输入数据类型来判断如何处理:string,address,uint256
					paramName := fmt.Sprintf("param%d", num)
					if vv.Type == "address" {
						func_nogas_data.InputParams += " ," + paramName
					} else if vv.Type == "uint256" || vv.Type == "uint" {
						func_nogas_data.InputParams += fmt.Sprintf(",%s", paramName)
					} else if vv.Type == "string" {
						func_nogas_data.InputParams += " ," + paramName
					} else {
						func_nogas_data.InputParams += " ," + paramName
					}

				}
				func_nogas_data.InputParams += ")"
				//输入参数
				num := 0
				strOutPuts := ""
				for _, _ = range v.Outputs {
					strOutPuts = fmt.Sprintf("%sdata%d,", strOutPuts, num)
					num++
				}
				strOutPuts += "err"
				func_nogas_data.OutParams = strOutPuts

				//模版的执行
				err = nogas_temp.Execute(outfile, &func_nogas_data)
				if err != nil {
					fmt.Println("failed to template nogas Execute ", err)
					return err
				}
			} else {
				//需要消耗gas
				func_gas_data.FuncName = v.Name
				func_gas_data.PrepareParams = prepareData
				func_gas_data.InputParams = "(testClient.GetTransactOpts()"

				for num, vv := range v.Inputs {
					//需要根据输入数据类型来判断如何处理:string,address,uint256
					paramName := fmt.Sprintf("param%d", num)
					if vv.Type == "address" {
						func_gas_data.InputParams += " ," + paramName
					} else if vv.Type == "uint256" {
						func_gas_data.InputParams += fmt.Sprintf(" ,%s", paramName)
					} else if vv.Type == "string" {
						func_gas_data.InputParams += " ," + paramName
					} else {
						func_gas_data.InputParams += " ," + paramName
					}

				}
				func_gas_data.InputParams += ")"
				//模版的执行
				err = hasgas_temp.Execute(outfile, &func_gas_data)
				if err != nil {
					fmt.Println("failed to template hasgas Execute ", err)
					return err
				}
			}
		}
	}

	return nil
}

func Impl_main_code(runCodePath, basicName string) error {
	//1. 写到哪
	outfile, err := os.OpenFile("main.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile.Close()
	// 读取abi文件信息
	abiInfos, err := readAbi(runCodePath + "/" + basicName + ".abi")
	if err != nil {
		fmt.Println("failed to read abi", err)
		return err
	}
	funcNames := ""
	//"abc","def","
	num := 0
	var funcInfos []FuncInfo
	var funcInfo FuncInfo
	// 2- 第一个函数
	for _, v := range abiInfos {
		if v.Type != "constructor" {
			if num == 0 {
				//第一个
				funcNames += fmt.Sprintf(`"%s"`, v.Name)
			} else {
				funcNames += fmt.Sprintf(`,"%s"`, v.Name)
			}
			num++
			funcInfo.FuncName = strings.Title(v.Name)
			funcInfo.Num = num + 1
			funcInfos = append(funcInfos, funcInfo)
		}
	}
	main_str1 := fmt.Sprintf(test_run_main_temp, funcNames)
	_, err = outfile.WriteString(main_str1)
	if err != nil {
		fmt.Println("failed to write to main.go", err)
		return err
	}

	//建立一个模版，输出内容
	main_temp, err := template.New("main").Parse(test_build_main_temp)
	if err != nil {
		fmt.Println("failed to template main", err)
		return err
	}
	err = main_temp.Execute(outfile, funcInfos)
	if err != nil {
		fmt.Println("failed to Execute main", err)
		return err
	}
	return err
}

func Run(buildPath, contractName, runCodeName string) {
	basicName := strings.Replace(contractName, ".sol", "", -1)
	err := Impl_run_code(buildPath, runCodeName, basicName)
	if err != nil {
		log.Panic("failed to Impl_run_code:", err)
	}
	err = Impl_main_code(buildPath, basicName)
	if err != nil {
		log.Panic("failed to Impl_main_code:", err)
	}
}

func ShowAbi(abifile string) {
	infos, err := readAbi(abifile)
	if err != nil {
		log.Panic("failed to readAbi:", err)
	}
	for _, v := range infos {
		fmt.Printf("%+v\n", v)
	}
}

func makePrepareParams(params []InputsOutPuts) string {

	prepareData := ""

	for num, v := range params {
		if v.Type == "uint256" {
			prepareData += fmt.Sprintf("\ttemp%d, _ := strconv.Atoi(params[%d])\n", num, num)
			prepareData += fmt.Sprintf("\tparam%d := big.NewInt(int64(temp%d))\n", num, num)
		} else if strings.Contains(v.Type, "int") {
			prepareData += fmt.Sprintf("\tparam%d, _ := strconv.Atoi(params[%d])\n", num, num)
		} else {
			prepareData += fmt.Sprintf("\tparam%d := params[%d]\n", num, num)
		}
	}
	return prepareData
}
