package BLC

import "fmt"

func (cli *CLI) createWallet() {
	wallets,_ := NewWallets()
	wallets.CreateNewWallet()

	//把所有数据保存起来
	wallets.SaveWalletsToFile()

	fmt.Println(wallets.WalletsMap)
}


