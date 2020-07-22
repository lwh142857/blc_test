package BLC

import (

	//	"bolt"
	//	"log"
	//	"os"

	//	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

//表名称
const accountsTableName = "accounts"

type Account struct {
	Name    string
	Balance int
	Address string
	Days    int
	Hash    []byte
}

var AccountsPool []Account   //账户池
var P_AccountsPool []Account //概率账户池
var V_AccountsPool []Account //竞选节点池
var S_AccountsPool []Account //超级节点

//币龄随时间增加
func Increasecoinage() {
	for _, v := range AccountsPool {
		v.Days++
	}

}

//添加账户
func NewAccount(name string, balance int) Account {
	account := Account{
		Name:    name,
		Balance: balance,
		Address: name,
		Days:    0,
		Hash:    nil,
	}
	return account
}

//挖矿成功奖励
func MiningReward(address string) {
	for _, v := range AccountsPool {
		if v.Address == address {
			v.Balance += miningreward
			v.Days = 0
		}
	}
}
func (account *Account) Serialize() []byte {
	var buffer bytes.Buffer
	//新建编码对象
	encoder := gob.NewEncoder(&buffer)
	//编码（序列化）
	if err := encoder.Encode(account); err != nil {
		log.Panicf("serialized the account to []byte failed %v\n", err)
	}
	return buffer.Bytes()
}

//数据反序列化
func DeserializeAccount(blockBytes []byte) *Account {
	var account Account
	//新建decoder对象
	//fmt.Printf("blockBytes : %v\n",blockBytes)      -----------调试bug
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&account); err != nil {
		log.Panicf("deserialized []byte to account failed %v\n", err)
	}
	return &Account{}
}
