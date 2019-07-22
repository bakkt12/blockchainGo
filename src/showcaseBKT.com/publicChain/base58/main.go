package main

import (
	"fmt"
	"showcaseBKT.com/publicChain/base58/BLC"
)
/**
 地址就用 ：URLEncoding
字符用：StdEncoding
 */

func main() {
	//msg := "Hello, 世界"
	//1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
	msg:= []byte("Hello, 世界OIL +++__----++/;/;///");
	encoded := BLC.Base58Encode(msg)
	//encoded := base64.StdEncoding.EncodeToString([]byte(msg))

	//fmt.Println(encoded)
	fmt.Printf("%s \n",encoded)
	fmt.Printf("%x \n",encoded)

	decoded := BLC.Base58Decode(encoded)

	fmt.Printf("原来 %s \n:",msg)
	fmt.Println(string(decoded))
}
