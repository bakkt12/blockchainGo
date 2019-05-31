package main

import (
	"fmt"
	"math/big"
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

	var hashInt big.Int
	blockchainIterator = blockchain.Iterator()
	for i := 0; i < 500; i++ {
		fmt.Printf("blockchain :%x \n", blockchainIterator.CurrentHash)
		blockchainIterator = blockchainIterator.Next()

		hashInt.SetBytes(blockchainIterator.CurrentHash)

		if( hashInt .Cmp(big.NewInt(0))==0) {
			fmt.Printf("%x \n",blockchainIterator.CurrentHash)
			break;
		}
	}
}

func  cli()  {
	blockchain:= BLC.NewBlockchain()

	cli:= BLC.CLI{blockchain}

	cli.Run()
}
func main() {
	//cli()

	var a [][]string

	a = append(a,[]string{"sss"})

	fmt.Printf("%s",a)
}


