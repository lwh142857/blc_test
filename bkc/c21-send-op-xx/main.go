package main

import "bkc/c21-send-op-xx/BLC"

//启动
func main(){
	/*bc := BLC.CreateBlockChainWithGenesisBlock()
	bc.AddBlock([]byte("a send 100 eth to b"))
	bc.AddBlock([]byte("b send 100 eth to c"))
	bc.AddBlock([]byte("c send 100 eth to d"))

	bc.PrintChain()*/

	cli :=BLC.CLI{}
	cli.Run()

}