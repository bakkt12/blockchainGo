package BLC

import "fmt"

//返回所有钱包地址
func (cli *CLI) getaddresslists() {
	//打印出所有
	fmt.Println("打印所有钱包地址..")

	wallets := NewWallets()
	for key, _ := range wallets.WalletsMap {
		fmt.Println(key)
	}
}
