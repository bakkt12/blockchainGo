package main

import (
	"fmt"
	"showcaseBKT.com/publicChain/Serialize-DeserializeBlock/BLC"
)

func AddBlock() {

	blockchain := BLC.NewBlockchain()
	fmt.Printf("NewBlockchain tip: %x\n", blockchain.Tip)
	blockchain.AddBlock("send ben0 btc")
	blockchain.AddBlock("send ben1 btc")
	blockchain.AddBlock("send ben2 btc")
	blockchain.AddBlock("send ben3 btc")
	blockchain.AddBlock("send ben4 btc")
	blockchain.AddBlock("send ben5 btc")
	blockchain.AddBlock("send ben6 btc")
	blockchain.AddBlock("send ben7 btc")
	blockchain.AddBlock("send ben8 btc")
	blockchain.AddBlock("send ben9 btc")
	blockchain.AddBlock("send ben10 btc")
	blockchain.AddBlock("send ben11 btc")
	blockchain.AddBlock("send ben12 btc")
	fmt.Println("============ end addblock===========")
}

func BlockIteraor() {
	blockchain := BLC.LastBlockchain()
	var blockchainIterator *BLC.BlockchainIterator
	fmt.Println("============ iterator===========")

	blockchainIterator = blockchain.Iterator()
	for i := 0; i < 80; i++ {
		fmt.Printf("blockchain :%x \n", blockchainIterator.CurrentHash)
		blockchainIterator = blockchainIterator.Next()
	}
}

func main() {
	//AddBlock()
	BlockIteraor()
}
