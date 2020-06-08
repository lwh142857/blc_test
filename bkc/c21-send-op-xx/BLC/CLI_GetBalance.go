package BLC

import "fmt"

//查询余额
func (cli *CLI) getBalance(from string){
	//查找改地址的UTXO
	//获取区块链对象
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()  //关闭实例对象
	amount:=blockchain.getbalance(from)
	fmt.Printf("\t地址 [%s] 的余额 ： [%d]\n",from,amount)
}