package main

import (
	"fmt"
	"kongyixueyuan.com/publicChain/part2-proof-of-work/BLC"
)

func main() {

	blockchain := BLC.NewBlockchain();
	fmt.Println(blockchain)
}
