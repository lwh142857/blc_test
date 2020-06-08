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
	target = target.Lsh(target,256-targetBit)
	return &ProofOfStake{block,target}
}

func (ProofOfStake *ProofOfStake) Run() ([]byte,int){
	var nonce=0
	var hashInt big.Int
	var hash [32]byte 
	dataBytes := proofOfWork.prepareData(int64(nonce))
	//生成hash
	hash = sha256.Sum256(dataBytes)  //sha256.Sum256返回的是[]byte
	//将hash存储到hashInt
	hashInt.SetBytes(hash[:]) 
	return hash[:],nonce
}

//随机得出挖矿地址（挖矿概率跟代币数量与币龄有关）
func getMineNodeAddress() string {
	bInt := big.NewInt(int64(len(Node.P_Nodespool)))
	//得出一个随机数，最大不超过随机节点池的大小
	rInt, err := rand.Int(rand.Reader, bInt)
	if err != nil {
		log.Panic(err)
	}
	Node.P_Nodespool(rInt.Int64()).Account.days = 0
	return Node.P_Nodespool[int(rInt.Int64())].address
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