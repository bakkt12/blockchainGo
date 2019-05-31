package BLC

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"time"
)

type CLI struct {
	BC *Blockchain
}

//直接打印usage信息
func (cli *CLI) printUsage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data BLOCK_DATA - add a block")
	fmt.Println("\tprintchain -print all the blocks of the blockchain")
}

//判断终端参数的个
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	var blockchainIterator *BlockchainIterator
	blockchainIterator = cli.BC.Iterator()
	var hashBigInt big.Int
	fmt.Println("")
	for {

		err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			//通过hash获取到区块字节数组
			currentBlockBytes := b.Get([]byte(blockchainIterator.CurrentHash))
			currentBlock := DeserializeBlock(currentBlockBytes)
			fmt.Printf("Data:%s \n", currentBlock.Data)
			fmt.Printf("PrevBlockHash:%x \n", currentBlock.PrevBlockHash)
			fmt.Printf("Hash:%x \n", currentBlock.Hash)
			fmt.Printf("Nonce:%d \n", currentBlock.Nonce)
			fmt.Printf("Timestamp:%s \n\n", time.Unix(currentBlock.Timestamp, 0).Format("2006-01-02 15:04:05"))

			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		//获取下一个迭代器
		blockchainIterator = blockchainIterator.Next()

		//是否到达创世区块
		hashBigInt.SetBytes(blockchainIterator.CurrentHash)
		if (hashBigInt.Cmp(big.NewInt(0)) == 0) {
			break;
		}

	}
}

func (cli *CLI) addBlock(data string) {
	cli.BC.AddBlock(data)
}
func (cli *CLI) Run() {
	//判断终端参数的个数 如果没有参数直接打印usage信息
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block")

	switch os.Args[1] {

	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsage()
			os.Exit(1)
		}
		//fmt.Printf("Data:" + *addBlockData)
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}
