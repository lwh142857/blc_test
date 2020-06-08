##1.转账逻辑完整与UTXO查找优化<br>
1.查找可用UTXO的函数FindSpendableTTXO()<br>
2.实现通过UTXO查询进行转账，修改NewSimpleTransaction()<br>
3.添加ProofOfStake<br>



go build -o bc.exe main.go<br>
bc.exe getbalance --address<br>
bc.exe printchain   <br>
bc.exe createblockchain<br>
bc.exe send --from "[\"troytan\",\"Alice\"]" --to "[\"Alice\",\"troytan\"]" --amount "[\"5\",\"2\"]"<br>


