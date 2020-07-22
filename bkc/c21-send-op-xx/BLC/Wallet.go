package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//钱包管理
//校验和长度
const addressCheckSumLen = 4

//钱包基本结构
type Wallet struct {
	//1.私钥
	PrivateKey ecdsa.PrivateKey
	//2.公钥
	PublicKey []byte
}

//创建一个钱包
func NewWallet() *Wallet {
	//公钥-私钥赋值
	privatekKey, publicKey := newKeyPair()
	return &Wallet{PrivateKey: privatekKey, PublicKey: publicKey}
}

//通过钱包生成公钥-私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//1.获取一个椭圆
	curve := elliptic.P256()
	//2.通过椭圆生成私钥
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if nil != err {
		log.Panicf("ecdsa generate private key failed! %v\n", err)
	}
	//3.通过私钥生成公钥
	pubKey := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return *priv, pubKey
}

//生成地址
//实现双哈希
func Ripemd160Hash(pubkey []byte) []byte {
	//sha256
	hash256 := sha256.New()
	hash256.Write(pubkey)
	hash := hash256.Sum(nil)
	//ripemd160
	rmd160 := ripemd160.New()
	rmd160.Write(hash)
	return rmd160.Sum(nil)
}

//生成校验和
func CheckSum(input []byte) []byte {
	first_hash := sha256.Sum256(input)
	second_hash := sha256.Sum256(first_hash[:])
	return second_hash[:addressCheckSumLen]
}

//通过钱包（公钥）获取地址
func (w *Wallet) GetAddress() []byte {
	//1.获取hash160
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	//2.获取校验和
	checkSumBytes := CheckSum(ripemd160Hash)
	//3.拼接
	addressBytes := append(ripemd160Hash, checkSumBytes...)
	//4.base58编码
	b58Bytes := Base58Encode(addressBytes)
	return b58Bytes
}

//判断地址的有效性
func IsValidForAddress(addressBytes []byte) bool {
	//1.地址通过base58Decode进行解码
	pubkey_checkSumByte := Base58Decode(addressBytes)
	//2.拆分，进行校验和校验
	checkSumBytes := pubkey_checkSumByte[len(pubkey_checkSumByte)-addressCheckSumLen:]
	ripemd160Hash := pubkey_checkSumByte[:len(pubkey_checkSumByte)-addressCheckSumLen]
	//3.生成
	checkBytes := CheckSum(ripemd160Hash)
	//4.比较
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}
