package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) PrintChain() {
	//检查是否有数据存在
	existDB := DBExists()
	if existDB == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	blockchian := BlockchainObject();
	defer blockchian.DB.Close()
	blockchian.Printchain()
}

