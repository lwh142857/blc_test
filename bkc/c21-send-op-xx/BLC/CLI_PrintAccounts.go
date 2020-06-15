package BLC

import (
	"fmt"
)

//查询余额
func (cli *CLI) printaccounts() {
	//查找改地址的UTXO
	//获取区块链对象
	blockchain := BlockchainObject()
	defer blockchain.DB.Close() //关闭实例对象

	for i := 0; i < len(AccountsPool); i++ {
		amount := blockchain.getbalance(AccountsPool[i].Address)
		fmt.Printf("\t账户名 [%s] \n地址 [%s] \n余额 ： [%d]\n", AccountsPool[i].Name, AccountsPool[i].Address, amount)
	}

}
