package BLC

//交易输出管理
//输出结构
type TxOutput struct {
	//金额
	Value int  //大写才能导出金额
	//用户名
	ScriptPubkey string
}

//验证当前UTXo是否属于指定的地址
func (txOutput *TxOutput) CheckPubkeyWithAddress (address string) bool{
	return address == txOutput.ScriptPubkey
}