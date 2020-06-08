package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"os"
)

//参数数量的检测函数
func IsValidArgs(){
	if len(os.Args) < 2{
		PrintUsage()
		//直接退出
		os.Exit(1)
	}
}

//实现int64转[]byte
func IntToHex(data int64) []byte{
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer,binary.BigEndian,data)
	if err != nil{
		log.Panicf("int transact to []byte failed! %v\n",err)
	}
	return buffer.Bytes()
}

//标准json格式转切片
func JSONToSlice(jsonString string)[]string{
	var strSlice []string
	//json
	if err:=json.Unmarshal([]byte(jsonString),&strSlice);err!=nil{
		log.Panicf("json to []string failed! %v\n",err)
	}
	return strSlice
}