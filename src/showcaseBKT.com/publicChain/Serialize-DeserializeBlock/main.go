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
	blockchain := BLC.NewBlockchain()
	cli := BLC.CLI{blockchain}
	cli.Run()
	//spentTXOs := make(map[string][]string)
	//for k, v := range spentTXOs {
	//	fmt.Printf("index %s,%d, \n", k, v)
	//}
	//nubmers := [] int{1, 2, 3, 4, 5, 6, 7, 8, 10, 1, 2, 3, 4, 5,6,6}
	//for index, x := range nubmers {
	//	vString :=strconv.Itoa(x)
	//	fmt.Printf("第 %d 位 x 的值 = %d ,v: %s\n", index,x ,vString )
	//	spentTXOs[vString ] = append(spentTXOs[vString ], vString )
	//}
	//
	//for k, v := range spentTXOs {
	//	fmt.Printf("index %s,%s, \n", k, v)
	//}
}
