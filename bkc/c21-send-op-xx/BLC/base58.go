package BLC

import (
	"bytes"
	"math/big"
)

//base58编码实现
//1.生成一个base58的编码基数表

var b58Alphabet = []byte("" +
	"123456789" +
	"abcdefghijkmnopqrstuvwxyz" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ")

//编码函数
func Base58Encode(input []byte) []byte {
	var result []byte //代表编码结果
	//big.Int
	x := big.NewInt(0).SetBytes(input)
	//求余的基本长度
	base := big.NewInt(int64(len(b58Alphabet)))
	//求余和商

	zero := big.NewInt(0)
	//设置余数，代表base58基数表的索引位置
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		//得到的result是个倒序的base58编码结果
		result = append(result, b58Alphabet[mod.Int64()])
	}

	//反转切片
	Reverse(result)
	//添加一个前缀，代表是一个地址
	result = append([]byte{b58Alphabet[0]}, result...)
	return result
}

//反转切片函数
func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

//解码函数
//传进来的input是编码结果
/*
   1Jh83
*/
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 1
	//1.去掉前缀
	data := input[zeroBytes:]
	for _, b := range data {
		charIndex := bytes.IndexByte(b58Alphabet, b) //内部函数，返回字符在切片中第一次出现的索引
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()

	return decoded
}
