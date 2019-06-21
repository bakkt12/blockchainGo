package BLC

import (
	"fmt"
	"os"
)



//创建创世区块
func (cli *CLI) createGenesisBlockchain(genesis string) {
	if DBExists() {
		fmt.Println("创世区块已经存在")
		os.Exit(1)
	}
	CreateBlockchainWithGenesisBlock(genesis)
}