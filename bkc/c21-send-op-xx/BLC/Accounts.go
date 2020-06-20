package BLC

type Account struct {
	Name    string
	Balance int
	Address string
	Days    int
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
func AddNewAccount(name string, balance int) Account {
	account := Account{
		Name:    name,
		Balance: balance,
		Address: name,
		Days:    0,
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
