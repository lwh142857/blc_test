package BLC

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"
	//	"crypto/sha256"
	//	"fmt"
)

//POS生成区块

type ProofOfStake struct {
	//需要共识的区块
	Block *Block
	//目标难度的哈希
	target *big.Int //大数据存储
}

//创建一个POS对象
func NewProofOfStake(block *Block) *ProofOfStake {
	//1.创建一个初始值为1的target
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfStake{block, target}
}

//执行函数
func (ProofOfStake *ProofOfStake) Run() ([]byte, int) {
	var nonce = 0
	var hashInt big.Int
	var hash [32]byte
	dataBytes := ProofOfStake.prepareData(int64(nonce))
	//生成hash
	hash = sha256.Sum256(dataBytes) //sha256.Sum256返回的是[]byte
	//将hash存储到hashInt
	hashInt.SetBytes(hash[:])
	return hash[:], nonce
}

//随机得出挖矿地址（挖矿概率跟代币数量与币龄有关）
func getMineNodeAddress() string {
	bInt := big.NewInt(int64(len(P_AccountsPool)))
	//得出一个随机数，最大不超过随机节点池的大小
	rInt, err := rand.Int(rand.Reader, bInt)
	if err != nil {
		log.Panic(err)
	}
	P_AccountsPool[int(rInt.Int64())].Days = 0
	return P_AccountsPool[int(rInt.Int64())].Address
}

//拼接区块属性，进行哈希计算
func (pos *ProofOfStake) prepareData(nonce int64) []byte {

	timeStampBytes := IntToHex(pos.Block.TimeStamp)
	heightBytes := IntToHex(pos.Block.Heigth)
	data := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		pos.Block.PrevBlockHash,
		pos.Block.HashTransaction(),
		IntToHex(targetBit),
		IntToHex(nonce),
	}, []byte{})
	return data
}

//初始化
func Pos_init() {
	AccountsPool = append(AccountsPool, AddNewAccount("1", 3200))
	AccountsPool = append(AccountsPool, AddNewAccount("2", 6400))
	//初始化随机节点池（挖矿概率与代币数量和币龄有关）
	for _, v := range AccountsPool {
		for i := 0; i <= v.Balance*v.Days; i++ {
			P_AccountsPool = append(P_AccountsPool, v)
		}
	}
}
