package BLC

import (
	"fmt"
	"os"
)

//发起交易
//send --from "[\"troytan\"]" --to "[\"Alice\"]" --amount "[\"3\"]"
/*区块链完整信息
------------------------
        Hash: 000040b0d92d45cd2703b1d5e44a714197a50fb3a186026fec2e08f100a79d06
        PrevBlockHash: 0000c49e02fa0f79b5aa97a76830ac5f7c56da401f774094b7bf8da45d25f272
        TimeStamp: 1588858740
        Heigth: 2
        Nonce: 26846
        Txs: [0xc000080050]
                tx-hash : df01942e6b4f195dd9a8596331182d5fcc2faa6263f02700a2f25d66f3f7696b
                输入...
                        vin-txHash : 36663932653030333737623931646462313636306665323232656237316432346361636564353661373337613938366266383037653663643838356532326637
                        vin-vout : 0
                        vin-scriptSig : troytan
                输出...
                        vout-value : 3
                        vout-scriptPubkey : Alice
                        vout-value : 7
                        vout-scriptPubkey : troytan
------------------------
        Hash: 0000c49e02fa0f79b5aa97a76830ac5f7c56da401f774094b7bf8da45d25f272
        PrevBlockHash:
        TimeStamp: 1588850017
        Heigth: 1
        Nonce: 62522
        Txs: [0xc0000800a0]
                tx-hash : 6f92e00377b91ddb1660fe222eb71d24caced56a737a986bf807e6cd885e22f7
                输入...
                        vin-txHash :
                        vin-vout : -1
                        vin-scriptSig : system reward
                输出...
                        vout-value : 10
                        vout-scriptPubkey : troytan


*/

func (cli *CLI) send(from,to,amount []string){
	if !dbExit(){
		fmt.Println("数据库不存在....")
		os.Exit(1)
	}
	//获取区块链对象
	blockchain :=BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from,to,amount)
}