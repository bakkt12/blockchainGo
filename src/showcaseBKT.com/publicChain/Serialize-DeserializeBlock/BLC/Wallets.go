package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "Wallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

//创建钱包的集合
func NewWallets() (*Wallets, error) {
	wallets := &Wallets{}
	wallets.WalletsMap = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return wallets, err;
}

func (ws *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	fmt.Sprintf("create new wallet: %s \n", wallet.Getaddress())
	ws.WalletsMap[string(wallet.Getaddress())] = wallet
}

//获取所有的钱包信息
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}

//通过钱包地址获取钱包对象
func (ws *Wallets) GetWallet(address string) *Wallet {
	return ws.WalletsMap[address]
}

//加载钱包文件
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	//  注册目的为是的可以序列化任何类型 （接口..）
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsMap = wallets.WalletsMap
	return nil
}

//把钱包信息保存到一个文件中
func (ws *Wallets) SaveWalletsToFile() {
	var content bytes.Buffer
	//  注册目的为是的可以序列化任何类型 （接口..）
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&ws)
	if err != nil {
		log.Panic(err)
	}
	//将序列化后的数据写入到文件 原来的数据会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
