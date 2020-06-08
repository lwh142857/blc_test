##1.转账逻辑完整与UTXO查找优化<br>
1.查找可用UTXO的函数FindSpendableTTXO()<br>
2.实现通过UTXO查询进行转账，修改NewSimpleTransaction()<br>
3.添加ProofOfStake和Accounts<br>


操作命令<br>

go build -o bc.exe main.go          初始化<br>
bc.exe getbalance --address         查询余额<br>
bc.exe printchain                   打印区块链<br>
bc.exe createblockchain             创建区块链<br>
bc.exe send --from "[\"troytan\",\"Alice\"]" --to "[\"Alice\",\"troytan\"]" --amount "[\"5\",\"2\"]"            转账 <br>


