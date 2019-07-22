package main

import (
	"golang.org/x/crypto/ripemd160"
	"fmt"
)
func main() {
	//160位  /8 = 20字节
	hasher := ripemd160.New()
	hasher.Write([]byte("我最喜欢叶昊宁"))
	bytes := hasher.Sum(nil)
	fmt.Printf("%x \n", bytes)
}
