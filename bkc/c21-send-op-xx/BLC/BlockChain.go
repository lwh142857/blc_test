package BLC

import (
	"bolt"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
)

//数据库名称
const dbName = "block.db"

//表名称
const blockTableName = "blocks"

//区块链基本结构
type BlockChain struct {
	//Blocks []*Block //区块的切片
	DB  *bolt.DB //数据库对象
	Tip []byte   //保存最新区块哈希值

}

//账户结构
type AccountData struct {
	DB  *bolt.DB //数据库对象
	Num []byte
}

//判断数据库文件是否存在
func dbExit() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		//数据库文件不存在
		return false
	}
	return true
}

//初始化区块链
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {
	if dbExit() {
		//文件已存在，说明创世区块已存在
		fmt.Println("创世区块已存在...")
		os.Exit(1)
	}

	//保存最新区块的哈希值
	var blockHash []byte
	//1.创建或打开一个数据库
	//w写4 r读2 x执行1
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("create db[%s] failed %v\n", dbName, err)
	}
	//2.创建桶,把生成的创世区块存入数据库中
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		//区块链表
		if b == nil {
			//没找到桶
			b, err := tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panicf("create bucket[%s] failed %v\n", blockTableName, err)
			}
			//生成一个coinbase交易，它是包括：这条交易哈希、vin、vout
			txCoinbase := NewCoinbaseTransaction(address)

			//生成创世区块，块里是包括多条交易的
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			//存储
			//1.key,value分别以什么数据代表
			//2.如何把block结构存入数据库中---序列化
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize()) //key是hash，value是序列化的结果-----无论是key还是value都是[]byte
			if err != nil {
				log.Panicf("insert the genesis block failed %v\n", err)
			}
			blockHash = genesisBlock.Hash
			//存储最新区块的哈希
			err = b.Put([]byte("1"), genesisBlock.Hash)
			if err != nil {
				log.Panicf("save the hash of genesis block failed %v\n", err)
			}
		}

		if b != nil {
			blockHash = b.Get([]byte("1"))
		}
		return nil
	})
	return &BlockChain{db, blockHash}
}

//添加账户到表中
func (bc *AccountData) AddAccount(name string, balance int) {
	//更新区块数据（insert）
	bc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取数据库桶
		b := tx.Bucket([]byte(accountsTableName))
		if b != nil {
			fmt.Printf("latest hash: %v\n", b.Get([]byte("1")))
			//1.新建
			newaccount := NewAccount(name, balance)
			fmt.Printf("the hash of the account %x\n", newaccount.Hash)
			//2.存入数据库
			err := b.Put(newaccount.Hash, newaccount.Serialize())
			if err != nil {
				log.Panicf("insert the new account to db failed %v", err)
			}
			//更新最新区块的哈希（数据库中的）
		}
		return nil
	})
}

//添加区块到区块链中
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//更新区块数据（insert）
	bc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取数据库桶
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			fmt.Printf("latest hash: %v\n", b.Get([]byte("1")))
			//2.获取最后插入的区块
			blockBytes := b.Get(bc.Tip) //或者写成 blockBytes:=b.Get([]byte("1"))
			//fmt.Printf("tip--%v\n",bc.Tip)   -----------调试bug
			//3.区块数据的反序列化
			latest_block := DeserializeBlock(blockBytes)
			//4.新建一个区块
			newBlock := NewBlock(latest_block.Height+1, latest_block.Hash, txs)
			fmt.Printf("the hash of the block %x\n", newBlock.Hash)
			//5.存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panicf("insert the new block to db failed %v", err)
			}
			//更新最新区块的哈希（数据库中的）
			err = b.Put([]byte("1"), newBlock.Hash)
			if err != nil {
				log.Panicf("update the latest block hash to db failed %v", err)
			}
			//更新区块链对象中的最新区块哈希
			bc.Tip = newBlock.Hash
		}
		return nil
	})
}

//遍历数据库，输出所有区块信息
func (bc *BlockChain) PrintChain() {
	//读取数据库
	fmt.Println("区块链完整信息")
	var curBlock *Block
	bcit := bc.Interator() //获取迭代器对象
	//循环读取
	//退出条件
	for {
		fmt.Println("------------------------")
		curBlock = bcit.Next()
		fmt.Printf("\tHash: %x\n", curBlock.Hash)
		fmt.Printf("\tPrevBlockHash: %x\n", curBlock.PrevBlockHash)
		fmt.Printf("\tTimeStamp: %v\n", curBlock.TimeStamp)
		fmt.Printf("\tHeight: %d\n", curBlock.Height)
		fmt.Printf("\tNonce: %d\n", curBlock.Nonce)
		fmt.Printf("\tTxs: %v\n", curBlock.Txs)
		for _, tx := range curBlock.Txs {
			fmt.Printf("\t\ttx-hash : %x\n", tx.TxHash)
			fmt.Printf("\t\t输入...\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\t\tvin-txHash : %x\n", vin.TxHash)
				fmt.Printf("\t\t\tvin-vout : %v\n", vin.Vout)
				fmt.Printf("\t\t\tvin-scriptSig : %s\n", vin.ScriptSig)
			}
			fmt.Printf("\t\t输出...\n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\t\tvout-value : %d\n", vout.Value)
				fmt.Printf("\t\t\tvout-scriptPubkey : %s\n", vout.ScriptPubkey)
			}
		}
		//退出条件
		//转换为big.Int
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PrevBlockHash)
		//比较
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			//遍历到创世区块
			break
		}
	}
}

//获取一个blockchain对象
func BlockchainObject() *BlockChain {
	//获取DB
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("open the db [%s] failed! %v\n", dbName, err)
	}
	//获取Tip
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	if err != nil {
		log.Panicf("get the blockchain object failed! %v\n", err)
	}

	return &BlockChain{db, tip}
}

//实现挖矿功能
//通过接收交易、生成区块
func (blockchain *BlockChain) MineNewBlock(from, to, amount []string) {
	//搁置交易生成步骤
	var txs []*Transaction
	var block *Block
	//遍历交易的参与者
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index]) //将amount字符串转换成int类型
		//生成新的交易
		tx := NewSimpleTransaction(address, to[index], value, blockchain)
		//追加到txs交易列表中
		txs = append(txs, tx)
	}
	//从数据库中获取最新区块
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//获取最新区块哈希值
			hash := b.Get([]byte("1"))
			//获取最新区块
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	//通过数据库中最新的区块去生成新的区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	//持久化新生成的区块到数据库中
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panicf("update the new block to db failed! %v\n", err)
			}

			//更新最新区块哈希值
			err = b.Put([]byte("1"), block.Hash)
			if err != nil {
				log.Panicf("update the latest block hash to db failed! %v\n", err)
			}
			blockchain.Tip = block.Hash
		}
		return nil
	})
}

//获取指定地址所有已花费输出
func (blockchain *BlockChain) SpentOutPuts(address string) map[string][]int {
	//已花费输出缓存
	//这里的map：key是string类型；value是一个整形数组[]int
	spentTXOutputs := make(map[string][]int)
	//获取迭代器对象
	bcit := blockchain.Interator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			//排除coinbase交易
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.CheckPubkeyWithAddress(address) {
						key := hex.EncodeToString(in.TxHash) //字节数组变字符串
						//添加到已花费输出
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return spentTXOutputs
}

/*
遍历查找区块链数据库中的每个区块的每一个交易
查找每个交易中的每个输出
判断每个输出是否满足以下条件
1.属于传入的地址
2.属于未花费
  2.1首先遍历一次区块链数据库，将所有已花费的Outpu存入一个缓存
  2.2再次遍历区块链数据库，检查每一个Vout是否包含在前面
*/

//查找指定地址UTXO
func (blockchain *BlockChain) UnUTXOS(address string) []*UTXO {
	//1.遍历数据库，查找所有与address相关的交易
	//获取迭代器
	bcit := blockchain.Interator()
	var unUTXOS []*UTXO //当前地址的未花费输出列表
	//获取指定地址所有已花费输出
	//迭代，不断获取下一个区块
	for {
		block := bcit.Next()
		//遍历区块中的每笔交易
		spentTXOutpus := blockchain.SpentOutPuts(address)
		for _, tx := range block.Txs {
			//跳转
		work:
			for index, vout := range tx.Vouts {
				//index：当前输出在当前交易的索引位置
				//vout:当前输出
				if vout.CheckPubkeyWithAddress(address) {
					//当前vout属于传入地址
					if len(spentTXOutpus) != 0 {
						var isSpentOutput bool //默认false
						for txHash, indexArray := range spentTXOutpus {
							for _, i := range indexArray {
								//txhash：当前输出所引用的交易哈希
								//indexArray:哈希关联的vout索引列表
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									//txHash ==tx.TxHash,说明当前的交易tx至少已经有输出被其他交易作为输入的引用
									//index==i说明正好当前花费被消费
									//跳转到最外层循环，判断下一个vout
									isSpentOutput = true
									continue work
								}
							}
						}
						if isSpentOutput == false {
							utxo := &UTXO{tx.TxHash, index, vout}
							unUTXOS = append(unUTXOS, utxo)
						}
					} else {
						//将当前地址所有输出都添加到未花费输出中
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unUTXOS
}

//查询余额
func (blockchain *BlockChain) getbalance(address string) int {
	var amount int //余额
	utxos := blockchain.UnUTXOS(address)
	for _, utxo := range utxos {
		amount += utxo.Ouput.Value
	}
	return amount
}

//查找指定地址可用UTXO,,超过amount就中断查找
//更新当前数据库中指定地址的UTXO数量
func (blockchain *BlockChain) FindSpendableUTXO(from string, amount int) (int, map[string][]int) {
	//可用的UTXO
	spendableUTXO := make(map[string][]int)

	var value int
	utxos := blockchain.UnUTXOS(from)
	//遍历UTXO
	for _, utxo := range utxos {
		value += utxo.Ouput.Value
		//计算交易哈希
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= amount {
			break
		}
	}
	//所有的都遍历完成，仍然小于amount，则说明资金不足
	if value < amount {
		fmt.Printf("地址 [%s] 金额不足,当前余额[%d],转账金额[%d]\n", from, value, amount)
	}
	return value, spendableUTXO
}
