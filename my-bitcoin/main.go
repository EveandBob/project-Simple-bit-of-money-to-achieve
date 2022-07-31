package main

import "bitcoin/my-bitcoin/BLC"

func main() {
	//创世区块
	Blockchain := BLC.CreateBlockChainWithGenesisBlock()
	//err := Blockchain.DB.Close()
	//if err != nil {
	//return
	//}
	Blockchain.AddBlockToBlockchain("send 2 bitcoin to satoshi")
	//fmt.Println("\n2 ok")
	Blockchain.AddBlockToBlockchain("send 3 bitcoin to satoshi")
	//fmt.Println("\n3 ok")
	//Blockchain.AddBlockToBlockchain("send 4 bitcoin to satoshi")
	//fmt.Println("\n4 ok")
	Blockchain.PrintChain()
}
