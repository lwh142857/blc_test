##1.转账逻辑完整与UTXO查找优化
1.查找可用UTXO的函数FindSpendableTTXO()
2.实现通过UTXO查询进行转账，修改NewSimpleTransaction()

go build -o bc.exe main.go /n
bc.exe getbalance --address /n
bc.exe printchain   <br>
bc.exe createblockchain
bc.exe send --from "[\"troytan\",\"Alice\"]" --to "[\"Alice\",\"troytan\"]" --amount "[\"5\",\"2\"]"