package main

import (
	"fmt"
	"kongyixueyuan.com/publicChain/part2-proof-of-work/BLC"
	"time"
)

const (
	date        = "2006-01-02"
	shortdate   = "06-01-02"
	times       = "15:04:02"
	shorttime   = "15:04"
	datetime    = "2006-01-02 15:04:02"
	newdatetime = "2006/01/02 15时04分02秒"
	newtime     = "15~04~02"
)

func main() {

	blockchain := BLC.NewBlockchain()
	blockchain.AddBlock("send ben 20 btc")
	blockchain.AddBlock("send ben -> bobi 50 btc")
	blockchain.AddBlock("send ben -> bobi 5 btc")
	blockchain.AddBlock("send ben -> bobi 6 btc")
	blockchain.AddBlock("send ben -> bobi 7 btc")
	blockchain.AddBlock("send ben -> bobi 8 btc")
	blockchain.AddBlock("send ben -> bobi 9 btc")

	fmt.Println("=======start block========")
	//	fmt.Println( blockchain)
	for _, block := range blockchain.Blocks {
		fmt.Printf("Data:%s ,PreBlockHash:%x ,Hash:%x,Timestamp:%v \n", block.Data, block.PrevBlockHash, block.Hash, time.Unix(block.Timestamp, 0).Format(newdatetime))
		fmt.Printf( "Hash  :%x \n",block.Hash)

		fmt.Printf( "Nonce :%b \n",block.Nonce)
	}
}
