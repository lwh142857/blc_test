package BLC

import (
	"fmt"
	"os"
)

//打印完整区块链信息
func (cli *CLI) printchain(){
	if !dbExit(){
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	//获取到blockchain对象实例
	blockchain.PrintChain()
}
