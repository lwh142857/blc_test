package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

//共识算法管理文件

//实现POW实例以及相关功能
//目标难度值

//0000 0000 0000 0000 1001 0001 0000 .... 0001 Hash256位
// 256位hash里面的前面至少有16个0
const targetBit = 16
type ProofOfWork struct {
	//需要共识的区块
	Block *Block
	//目标难度的哈希
	target *big.Int  //大数据存储
}

//创建一个POW对象
func NewProofOfWork(block *Block) *ProofOfWork{
	//1.创建一个初始值为1的target
	target :=big.NewInt(1)
	//难度 2，假设哈希是8位
	//0000 0001
	//8-2=6
	//0100 0000 64 左移6位，生成的哈希只要小于64就可以了

	//2.左移256-targetBit
	target = target.Lsh(target,256-targetBit)
	return &ProofOfWork{block,target}
}

//执行pow，比较哈希
//返回哈希值，以及碰撞的次数
func (proofOfWork *ProofOfWork) Run() ([]byte,int){
	//碰撞次数
	var nonce=0
	var hashInt big.Int
	var hash [32]byte //生成的哈希值,在外面先进行定义一下，从而当返回的时候，可以写return hash，具有可见性
	//无限循环，生成符合条件的哈希值
	for{
		//生成准备数据
		dataBytes := proofOfWork.prepareData(int64(nonce))
		//生成hash
		hash = sha256.Sum256(dataBytes)  //sha256.Sum256返回的是[]byte
		//将hash存储到hashInt
		hashInt.SetBytes(hash[:])     //sha256.Sum256是32个字节,是固定大小的字节数组，而block.hash是动态的数组
		//检测生成的哈希值是否符合条件
		//判断hashInt是否小于Block 里面的target
		//Cmp compares x and y and returns
		//-1 if x < y
		//0 if x==y
		//+1 if x > y
		//如下target是x，hashInt是y
		if proofOfWork.target.Cmp(&hashInt)==1{
			//找到了符合条件的哈希值，中断循环
			break
		}
		nonce++
	}
	fmt.Printf("碰撞次数：%d\n",nonce)
	return hash[:],nonce
}

func (pow *ProofOfWork)prepareData(nonce int64) []byte{
	//var data []byte
	//拼接区块属性，进行哈希计算
	timeStampBytes := IntToHex(pow.Block.TimeStamp)
	heightBytes := IntToHex(pow.Block.Heigth)
	data := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		pow.Block.PrevBlockHash,
		pow.Block.HashTransaction(),
		IntToHex(targetBit),
		IntToHex(nonce),
	},[]byte{})
	return data
}

