package main

import (
	"fmt"
	"kongyixueyuan.com/publicChain/Serialize-DeserializeBlock/BLC"
)

func main() {

	block :=BLC.Block{[]byte("send btc 3 to bakkt"),20000}
	fmt.Printf("%s \n",block.Data)
	fmt.Printf("%d \n",block.Nonce)


	bytes := block.Serialize()
	fmt.Println(bytes)

	blc:=BLC.DeserializeBlock(bytes)
	fmt.Printf("%s \n",blc.Data)
	fmt.Printf("%d \n",blc.Nonce)

}
