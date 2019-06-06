package main

import (
	"fmt"
	"math/big"
	"showcaseBKT.com/publicChain/Serialize-DeserializeBlock/BLC"
)

func BlockIteraor() {
	blockchain := BLC.LastBlockchain()
	var blockchainIterator *BLC.BlockchainIterator
	fmt.Println("============ iterator===========")

	var hashInt big.Int
	blockchainIterator = blockchain.Iterator()
	for i := 0; i < 500; i++ {
		fmt.Printf("blockchain :%x \n", blockchainIterator.CurrentHash)
		blockchainIterator = blockchainIterator.Next()

		hashInt.SetBytes(blockchainIterator.CurrentHash)

		if (hashInt.Cmp(big.NewInt(0)) == 0) {
			fmt.Printf("%x \n", blockchainIterator.CurrentHash)
			break;
		}
	}
}

//func cli() {
//	blockchain := BLC.NewBlockchain()
//	cli := BLC.CLI{blockchain}
//	cli.Run()
//}
//
//func printfBlock() {
//	blockchain := BLC.NewBlockchain()
//	cli := BLC.CLI{blockchain}
//	cli.PrintChain()
//}

func main() {
	cli := BLC.CLI{}
	cli.Run()
}
