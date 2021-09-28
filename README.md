## 整体介绍

chinaner是一个可自动化、基于fisco-bcos平台的智能合约测试工具，可自动化编译智能合约，生成部署智能合约和函数调用的测试工程项目（Go语言）。

## 依赖安装

### 1. 安装solc编译器

为了适应fisco-bcos平台，因此建议安装0.6.10版本的solc编译器。

官方项目地址： https://github.com/ethereum/solidity

通过github项目，可以找到对应的发行版本，以Linux系统为例，找到tag为0.6.10的版本即可，下载命令如下：

```sh
wget https://github.com/ethereum/solidity/releases/download/v0.6.10/solc-static-linux
chmod u+x solc-static-linux
mv solc-static-linux ~/bin/
ln -s ~/bin/solc-static-linux ~/bin/solc
```

`~bin`目录已经加入到了PATH目录。

测试solc命令是否安装成功。

```sh
$solc --version
solc, the solidity compiler commandline interface
Version: 0.6.10+commit.00c0fcaf.Linux.g++

```





### 2. 安装go-sdk的abigen工具

第二个需要安装的是FISCO-BCOS下的go-sdk工具，只需要下载源码包，直接编译即可。

```sh
git https://github.com/FISCO-BCOS/go-sdk.git
cd go-sdk/cmd/abigen
go build
```

通过上述命令就可以编译得到abigen工具，同样可以将其放到PATH环境变量的目录下。



## 使用演示



### 1. 下载编译chinaner

```sh
$ git clone https://github.com/yekai1003/chinaner.git

$ cd chinaner/

$ go build
$ ls
chinaner  compiler.go  go.mod  go.sum  LICENSE  main.go  README.md  templates
$ mv chinaner ~/bin/

```

编译完成后，可以将其放到PATH环境变量的路径下。



### 2. chinaner使用示例

查看帮助：

```sh
$ chinaner 

Version:
	1.0.0
Author:
	Yekai
Usage:
	chinaner build -file contract's name  -- for compiler contracts
	chinaner init -dir DIRNAME -sdkpath SDKPATH  -- init test dir
	chinaner generate -file contract's name  -- generate contract's code
	chinaner showabi -abi ABIFILE  -- showabi 


```

初始化一个工程目录，需要指定工程目录，以及fisco节点的SDK目录（证书相关文件所在目录）

```sh
$ chinaner init -dir yekaidemo -sdkpath ~/fisco/agencyA/sdk/
/usr/bin/mkdir -p yekaidemo/accounts
init cert files sucess!
/usr/bin/cp /home/jbf/fisco/agencyA/sdk//ca.crt yekaidemo
/usr/bin/cp /home/jbf/fisco/agencyA/sdk//sdk.crt yekaidemo
/usr/bin/cp /home/jbf/fisco/agencyA/sdk//sdk.key yekaidemo
done...

```

进入到工程目录下

```sh
cd yekaidemo
```



示例合约代码如下，将其保存为demo.sol（文件名要与合约名字一致）

```js
pragma solidity^0.6.10;


contract demo {
    string  mymsg;
    uint256 amount;
    mapping(bytes32=>bytes) hashDatas;
    
    constructor(string memory _msg, uint256 _amount) public {
        mymsg = _msg;
        amount = _amount;
    }
    function setMsg(string memory _msg) public {
        mymsg = _msg;
    }
    
    function getMsg() public view returns (string memory) {
        return mymsg;
    }
    
    function addHashData(bytes32 _hash, bytes memory _data) public  {
        require(hashDatas[_hash].length == 0, "hash already exists");
        hashDatas[_hash] = _data;
    }
    
    function getHashData(bytes32 _hash) public view returns (bytes memory) {
        return hashDatas[_hash];
    }
    
    function getTwoMsg() public view returns (string memory, uint256) {
        return (mymsg, amount);
    }
    
}
```

编译合约

```sh
$ chinaner build -file demo.sol 
/home/jbf/bin/solc --bin --abi -o ./contracts demo.sol
compile abi sucess!
/home/jbf/bin/abigen -sol demo.sol -pkg contracts -type demo -out contracts/demo.go
compile contracts to golang sucess!
done...

```

生成测试代码

```sh
$ chinaner generate -file demo.sol 
/usr/local/go/bin/go mod init yekaidemo
/usr/local/go/bin/go build
done...

```

查看目录

```sh
$ tree ./
./
├── accounts
│?? ├── 0x5946a2ec703e74ce91ac0703396be65daeb5ea99.pem
│?? └── 0x5946a2ec703e74ce91ac0703396be65daeb5ea99.public.pem
├── ca.crt
├── call.go
├── config.toml
├── contracts
│?? ├── demo.abi
│?? ├── demo.bin
│?? └── demo.go
├── demo.sol
├── go.mod
├── go.sum
├── main.go
├── sdk.crt
├── sdk.key
└── yekaidemo

```

### 3. 工程使用测试

查看帮助

```sh
$ ./yekaidemo 
Usage:
./yekaidemo deploy -- for deploy contract
./yekaidemo AddHashData -- for AddHashData 
./yekaidemo GetHashData -- for GetHashData 
./yekaidemo GetMsg -- for GetMsg 
./yekaidemo SetMsg -- for SetMsg 

```

部署合约，需要根据参数来决定传入几个参数

```sh
$ ./yekaidemo deploy zhangsan 100000
addr= 0x7C8C9a4Ee800a01C57cf0dd420f1DF5b8D46Ccb0 hash= 0xd400d5ae5ec965a1b9aec1424cf6a29e30556e204967d635efc2975d80682024

```

设置合约地址为环境变量CONTRACT_ADDR

```sh
export CONTRACT_ADDR=0x7C8C9a4Ee800a01C57cf0dd420f1DF5b8D46Ccb0
```

其他函数调用

```sh
# 获取两个返回值
$ ./yekaidemo GetTwoMsg
val0: zhangsan
val1: 100000
# 设置字符串数据
$ ./yekaidemo SetMsg lisi
0x43d56852d539efe516b9d2783059e5affdc9e820e72d9a615e7228b2085602f0
# 查看字符串数据
$ ./yekaidemo GetMsg
val0: lisi
# 尝试hash数据修改
$ ./yekaidemo AddHashData 0x43d56852d539efe516b9d2783059e5affdc9e820e72d9a615e7228b2085602f0 chinaner
0xfbc5329e1505f935e179dd11260b52daf8e37eee4c1552c968b9c9d5512d397a
# 获取hash存储的数据
$ ./yekaidemo GetHashData 0x43d56852d539efe516b9d2783059e5affdc9e820e72d9a615e7228b2085602f0 
val0: chinaner

```

## 说明

Solidity数据类型较多，对于一些特殊类型可能会存在不支持的情况，欢迎提issue一起讨论。