package main

import (
	"encoding/json"
	"showcaseBKT.com/publicChain/Serialize-DeserializeBlock/BLC"
)

//func BlockIteraor() {
//	blockchain := BLC.LastBlockchain()
//	var blockchainIterator *BLC.BlockchainIterator
//	fmt.Println("============ iterator===========")
//
//	var hashInt big.Int
//	blockchainIterator = blockchain.Iterator()
//	for i := 0; i < 500; i++ {
//		fmt.Printf("blockchain :%x \n", blockchainIterator.CurrentHash)
//		blockchainIterator = blockchainIterator.Next()
//
//		hashInt.SetBytes(blockchainIterator.CurrentHash)
//
//		if (hashInt.Cmp(big.NewInt(0)) == 0) {
//			fmt.Printf("%x \n", blockchainIterator.CurrentHash)
//			break;
//		}
//	}
//}

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

type StuRead struct {
	Name  interface{}
	Age   interface{}
	HIgh  interface{}
	Class json.RawMessage `json:"class"` //注意这里
}

type Class struct {
	Name  string
	Grade int
}

func main() {

	//main send -from [\"helloggggggggg\",\"apple\",\"java\",\"golang\"] -to [\"bbc\",\"btc\",\"bkc\",\"blc\"] -amount [\"10\",\"20\"]

	cli := BLC.CLI{}
	cli.Run()

}
