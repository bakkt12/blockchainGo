package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	//256位  /8 = 32字节
	hasher := sha256.New()
	hasher.Write([]byte("我最喜欢叶昊宁"))
	bytes := hasher.Sum(nil)
	fmt.Printf("%d \n", bytes)
}
