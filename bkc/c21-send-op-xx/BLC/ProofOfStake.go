package BLC

import(
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)
//POS生成区块

const targetBit = 16
type ProofOfStake struct {
	//需要共识的区块
	Block *Block
	//目标难度的哈希
	target *big.Int  //大数据存储
}

//创建一个POS对象
func NewProofOfStake(block *Block) *ProofOfStake{
	//1.创建一个初始值为1的target
	target :=big.NewInt(1)
	//难度 2，假设哈希是8位
	//0000 0001
	//8-2=6
	//0100 0000 64 左移6位，生成的哈希只要小于64就可以了

	//2.左移256-targetBit
	target = target.Lsh(target,256-targetBit)
	return &ProofOfStake{block,target}
}

func (ProofOfStake *ProofOfStake) Run() ([]byte,int){
	//碰撞次数
	var nonce=0
	var hashInt big.Int
	var hash [32]byte //生成的哈希值,在外面先进行定义一下，从而当返回的时候，可以写return hash，具有可见性
	//无限循环，生成符合条件的哈希值

	return hash[:],nonce
}


	//拼接区块属性，进行哈希计算
func (pos *ProofOfStake)prepareData(nonce int64) []byte{

	timeStampBytes := IntToHex(pos.Block.TimeStamp)
	heightBytes := IntToHex(pos.Block.Heigth)
	data := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		pos.Block.PrevBlockHash,
		pos.Block.HashTransaction(),
		IntToHex(targetBit),
		IntToHex(nonce),
	},[]byte{})
	return data
}