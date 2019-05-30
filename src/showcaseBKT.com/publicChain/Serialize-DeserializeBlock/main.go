package main

import (
	"fmt"
	"showcaseBKT.com/publicChain/Serialize-DeserializeBlock/BLC"
)

func main() {
  blockchain:= BLC.NewBlockchain()

  fmt.Println(blockchain)
  fmt.Printf("NewBlockchain tip: %x\n",blockchain.Tip)

	blockchain.AddBlock("send ben100 btc")
	fmt.Printf("tip: %x\n",blockchain.Tip)

	blockchain.AddBlock("send ben200 btc")
	fmt.Printf("tip: %x\n",blockchain.Tip)

	blockchain.AddBlock("send ben300 btc")
	fmt.Printf("tip: %x\n",blockchain.Tip)

	blockchain.AddBlock("send ben400 btc")
	fmt.Printf("tip: %x\n",blockchain.Tip)
}