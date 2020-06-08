package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

//对blockchain的命令行操作进行管理

//go build -o bc.exe main.go
//bc.exe
//Usage:
//        createblockchain -- 创建区块链
//        addblock --data DATA -- 添加区块
//        printchain -- 输出区块链信息

//bc.exe createblockchain
//碰撞次数：45315

//bc.exe printchain
//区块链完整信息
//------------------------
//        Hash: 0000d52efe1cc6df26eaa88c86e9abee2cdb2e66f16bc55f33932e972e55a5f0
//        PrevBlockHash:
//        TimeStamp: 1588736847
//        Data: [105 110 105 116 32 98 108 111 99 107 99 104 97 105 110]
//        Heigth: 1
//        Nonce: 45315

//bc.exe addblock --data "troytan send 100 eros to bob"
//latest hash: [0 0 213 46 254 28 198 223 38 234 168 140 134 233 171 238 44 219 46 102 241 107 197 95 51 147 46 151 46 85 165 240]
//碰撞次数：80759
//the hash of the block 0000156eb2a3da7e47be8b25f768e7644643d8cf610490248ed8c1981ff93cf6

//client对象
type CLI struct {
}

//用法展示
func PrintUsage(){
	fmt.Println("Usage:")
	//初始化区块链
	fmt.Printf("\tcreateblockchain --address address -- 创建区块链\n")
	//添加区块
	fmt.Printf("\taddblock --data DATA -- 添加区块\n")
	//打印完整的区块信息
	fmt.Printf("\tprintchain -- 输出区块链信息\n")
	//通过命令转账
	fmt.Printf("\t --from FROM --to TO --amount AMOUNT -- 发起转账\n")
	fmt.Printf("\t转账参数说明\n")
	fmt.Printf("\t\t--from FROM --转账原地址\n")
	fmt.Printf("\t\t--to TO --转账目的地址\n")
	fmt.Printf("\t\t--amount AMOUNT --转账金额\n")
	//查询余额
	fmt.Printf("\tgetbalance --address FROM --查询指定地址的余额\n")
	fmt.Printf("\t转账参数说明\n")
	fmt.Printf("\t\t--address --查询余额地址\n")
	//bc.exe send --from xiaxia --to sang --amount 1000
	//  FROM:[xiaxia]
	//  TO:[sang]
	//  Amount:[1000]
}


//添加区块
func (cli *CLI) addBlock(txs []*Transaction){
	//判断数据库是否存在
	if !dbExit(){
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	//获取到blockchain对象实例
	blockchain.AddBlock(txs)

}

//命令行运行函数
func (cli *CLI) Run(){
	//检测参数数量
	IsValidArgs()  //如果没有输入任何的命令，那么len(os.Args)==1

	//新建相关命令
	//添加区块
	addBlockCmd :=flag.NewFlagSet("addblock",flag.ExitOnError)
	//输出区块链完整信息
	printChainCmd :=flag.NewFlagSet("printchain",flag.ExitOnError)
	//创建区块链
	createBLCWithGenesisBlockCmd:= flag.NewFlagSet("createblockchain",flag.ExitOnError)
	//发起交易
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
	//查询余额的命令
	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)

	//数据参数
	//添加区块 比如*flagAddBlockArg，就是输入的命令行中，data命令后输入的字符串数据具体是什么
	flagAddBlockArg := addBlockCmd.String("data","send 100btc to player","添加区块数据")
	//区块链时指定的矿工地址（接收奖励）
	flagCreateBlockchainArg := createBLCWithGenesisBlockCmd.String("address","troytan","指定接收系统奖励的矿工地址")
	//发起交易参数
	flagSendFromArg := sendCmd.String("from","","转账原地址")
	flagSendToArg := sendCmd.String("to","","转账目标地址")
	flagSendAmountArg := sendCmd.String("amount","","转账金额")
	//查询余额命令行参数
	flagGetBalanceArg :=getBalanceCmd.String("address","","要查询的地址")

	//判断命令
	switch os.Args[1] {
	case "getbalance":
		//为什么下面的都是[2:]：因为没有输入任何命令时，len(os.Args)已经是一位了，再输入比如命令send 以后，判断send命令后的其他参数，那么显然是【0】，【1】，从2起
		if err:=getBalanceCmd.Parse(os.Args[2:]);err!=nil{
			log.Panicf("parse getBalanceCmd failed! %v\n",err)
		}
	case "send":
		if err:=sendCmd.Parse(os.Args[2:]);err!=nil{
			log.Panicf("parse sendCmd failed! %v\n",err)
		}
	case "addblock":
		if err:=addBlockCmd.Parse(os.Args[2:]);err!=nil{
			log.Panicf("parse addBlockCmd failed! %v\n",err)
		}
	case "printchain":
		if err:=printChainCmd.Parse(os.Args[2:]);err!=nil{
			log.Panicf("parse printchainCmd failed! %v\n",err)
		}
	case "createblockchain":
		if err:=createBLCWithGenesisBlockCmd.Parse(os.Args[2:]);err!=nil{
			log.Panicf("parse createBLCWithGenesisBlockCmd failed! %v\n",err)
		}
	default:
		//没有传递任何命令或者传递的命令不在上面的命令列表中
		PrintUsage()
	    os.Exit(1)
	}

	//查询余额，解析
	if getBalanceCmd.Parsed(){
		if *flagGetBalanceArg==""{
			fmt.Println("查询地址不能为空....")
			PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg)
	}

	//发起转账
	if sendCmd.Parsed(){
		if *flagSendFromArg ==""{
			fmt.Println("原地址不能为空....")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArg ==""{
			fmt.Println("目标地址不能为空....")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg ==""{
			fmt.Println("金额不能为空....")
			PrintUsage()
			os.Exit(1)
		}

		//bc.exe send --from "[\"troytan\"]" --to "[\"Alice\"]" --amount "[\"100\"]"
		//FROM:[[troytan]]
		//TO:[[Alice]]
		//Amount:[[100]]

		fmt.Printf("\tFROM:[%s]\n",JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n",JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAmount:[%s]\n",JSONToSlice(*flagSendAmountArg))
		cli.send(JSONToSlice(*flagSendFromArg),JSONToSlice(*flagSendToArg),JSONToSlice(*flagSendAmountArg))
	}
	//添加区块命令
	if addBlockCmd.Parsed(){
		if *flagAddBlockArg ==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}

	//输出区块信息
	if printChainCmd.Parsed(){
		cli.printchain()
	}
	//创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed(){
		if *flagCreateBlockchainArg ==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagCreateBlockchainArg)
	}

}