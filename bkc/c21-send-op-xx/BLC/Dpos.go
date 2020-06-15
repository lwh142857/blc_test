package BLC

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"
	"sort"
	//	"fmt"
)

const (
	voteNodeNum      = 100 //投票节点
	superNodeNum     = 10  //超级节点
	mineSuperNodeNum = 3
)

type Dpos struct {
	//需要共识的区块
	Block *Block
	//目标难度的哈希
	target *big.Int //大数据存储
}

//创建一个POS对象
func NewDPos(block *Block) *Dpos {
	//1.创建一个初始值为1的target
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &Dpos{block, target}
}

//投票
func voting() {
	for _, v := range AccountsPool {
		rInt, err := rand.Int(rand.Reader, big.NewInt(superNodeNum+1))
		if err != nil {
			log.Panic(err)
		}
		V_AccountsPool[int(rInt.Int64())].Balance += v.Balance
	}
}

//对挖矿节点进行排序
func sortMineNodes() {
	sort.Slice(V_AccountsPool, func(i, j int) bool {
		return V_AccountsPool[i].Balance > V_AccountsPool[j].Balance
	})
	S_AccountsPool = V_AccountsPool[:mineSuperNodeNum]
}

//执行函数
func (Dpos *Dpos) Run() []byte {
	var nonce = 0
	var hashInt big.Int
	var hash [32]byte
	dataBytes := Dpos.prepareData(int64(nonce))
	//生成hash
	hash = sha256.Sum256(dataBytes) //sha256.Sum256返回的是[]byte
	//将hash存储到hashInt
	hashInt.SetBytes(hash[:])
	//竞选+排序
	voting()
	sortMineNodes()

	return hash[:]
}

//拼接区块属性，进行哈希计算
func (dpos *Dpos) prepareData(nonce int64) []byte {

	timeStampBytes := IntToHex(dpos.Block.TimeStamp)
	heightBytes := IntToHex(dpos.Block.Heigth)
	data := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		dpos.Block.PrevBlockHash,
		dpos.Block.HashTransaction(),
		IntToHex(targetBit),
		IntToHex(nonce),
	}, []byte{})
	return data
}
