package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//区块基本结构和功能管理文件

//挖矿奖励
const miningreward = 5
const mod = "pow"

//实现一个最基本的区块结构
type Block struct {
	//注意这里面有三个参数不是byte类型，在拼接以计算哈希时，需要对int64类型进行转换

	TimeStamp     int64          //区块时间戳、代表区块时间
	Hash          []byte         //当前区块哈希
	PrevBlockHash []byte         //前区块哈希
	Heigth        int64          //区块高度
	Txs           []*Transaction //交易数据（交易列表）
	Nonce         int64          //在运行pow时生成的哈希变化值，也代表pow运行时动态修改的数据
	Mineraddress  string         //挖块的矿工地址
}

//构建区块
//创建新块时，时间戳在创建的时候生成，不同参数传递
//区块的哈希值需要经过pow共识算法的运算，计算出来，也不用外面进行参数传递
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block {
	block := Block{
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
		PrevBlockHash: prevBlockHash,
		Heigth:        height,
		Txs:           txs,
		Mineraddress:  nil,
		Nonce:         0,
	}
	//选择共识算法
	switch mod {
	case "pow":
		//通过POW生成新的哈希
		pow := NewProofOfWork(&block)
		//执行工作量证明算法
		hash, nonce := pow.Run()
		block.Hash = hash
		block.Nonce = int64(nonce)
	case "pos":
		//pos
		pos := NewProofOfStake(&block)
		block.Hash = pos.Run()
	case "dpos":
		//dpos
		dpos := NewDPos(&block)
		block.Hash = dpos.Run()
	}
	return &block
}

//生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

//区块结构序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	//新建编码对象
	encoder := gob.NewEncoder(&buffer)
	//编码（序列化）
	if err := encoder.Encode(block); err != nil {
		log.Panicf("serialized the block to []byte failed %v\n", err)
	}
	return buffer.Bytes()
}

//区块数据反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	//新建decoder对象
	//fmt.Printf("blockBytes : %v\n",blockBytes)      -----------调试bug
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&block); err != nil {
		log.Panicf("deserialized []byte to block failed %v\n", err)
	}
	return &block
}

//把指定区块中所有交易结构都序列化(类Merkle的哈希计算方法)
func (block *Block) HashTransaction() []byte {
	//每个tx.Hash都是一个字节数组[]byte,将所有交易拼接在一起，得到的就是二维数组
	//这个txHashes是对于某个区块而言里面所有交易放一起的哈希，有点Merkle根的意味
	var txHashes [][]byte
	//将指定区块中所有交易哈希进行拼接
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
