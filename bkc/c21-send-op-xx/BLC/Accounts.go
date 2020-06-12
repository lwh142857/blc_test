package BLC

import(
	"fmt"
	"blcok"
	"crypto/sha256"
	"time"
)

type Account struct{
	Name string
	Balance int
	Address string
	Days int
}
var AccountsPool[] Account   //账户池
var P_AccountsPool[] Account  //概率账户池
var V_AccountsPool[] Account //竞选节点池
var S_AccountsPool[] Account   //超级节点

//币龄随时间增加
func Increasecoinage()  {
	t := time.NewTicker(3*time.Second)
	for {	
		<-t.C:
		for _, v := range AccountsPool{
			AccountsPool[v].days++
			}
		}
}
//添加账户
func AddNewAccount(name string，balance float32)Account{
	account := Account{
		Name:name,
		Balance:balance,
		Address:strconv.Itoa(name + strconv.Itoa(balance)),
		Days:0,
	}
	return account
}