package main

import (
	"fmt"
	"kongyixueyuan.com/publicChain/Serialize-DeserializeBlock/BLC"
)

func main() {
  blockchain:= BLC.NewBlockchain()

  fmt.Println(blockchain)
  fmt.Printf("tip: %x\n",blockchain.Tip)

	blockchain.AddBlock("send ben100 btc")
	fmt.Printf("tip: %x\n",blockchain.Tip)
}