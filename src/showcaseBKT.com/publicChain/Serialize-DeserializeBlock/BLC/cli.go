package BLC

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
)

type CLI struct {
	BC *Blockchain
}

//直接打印usage信息
func (cli *CLI) printUsage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tcreate -address create blockchain.")
	fmt.Println("\tget    -address create blockchain.")
	fmt.Println("\tprintf -print all the blocks of the blockchain.")
}

//判断终端参数的个
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) PrintChain() {
	fmt.Println("============PrintChain==============================")

	//检查是否有数据存在
	if dbExists() == false {
		cli.printUsage()
		return
	}
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
			//fmt.Printf("Data:%s \n", currentBlock.Transcation)
			fmt.Println("============START==============================")
			fmt.Printf("PrevBlockHash:%x \n", currentBlock.PrevBlockHash)
			fmt.Printf("Hash			:%x \n", currentBlock.Hash)
			//fmt.Printf("Nonce			:%d \n", currentBlock.Nonce)
			//	fmt.Printf("Timestamp		:%s \n", time.Unix(currentBlock.Timestamp, 0).Format("2006-01-02 15:04:05"))
			for _, tx := range currentBlock.Transcation {
				fmt.Println("\t**************************")
				tx.printfTranscation()
				fmt.Println("\t**************************")
			}
			fmt.Println("===========END===============================")
			fmt.Println("")
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

func (cli *CLI) SendToken() {
	//1 .3->yjc
	//   3->bakkt

	//1.新建一个交易
	var noPackageTxs []*Transcation
	tx1 := NewUTXOTransaction("yhn", "BAKKT2", 5, cli.BC, noPackageTxs)
	noPackageTxs = append(noPackageTxs, tx1)

	tx2 := NewUTXOTransaction("BAKKT2", "YE11", 5, cli.BC, noPackageTxs)
	noPackageTxs = append(noPackageTxs, tx2)
	////
	//tx3 := NewUTXOTransaction("yhn", "LY11", 5, cli.BC, noPackageTxs)
	//noPackageTxs = append(noPackageTxs,tx3)

	cli.BC.MineBlock([]*Transcation{tx1, tx2 /*,tx2,tx3*/ })

}

func (cli *CLI) addBlock(data string) {
	cli.SendToken();
}
func (cli *CLI) Run() {
	//判断终端参数的个数 如果没有参数直接打印usage信息
	cli.validateArgs()

	crateBlockchainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	genenisAddress := crateBlockchainCmd.String("address", "", "create block and package dtata")
	printChainCmd := flag.NewFlagSet("printf", flag.ExitOnError)

	getBalanceCmd := flag.NewFlagSet("get", flag.ExitOnError)
	balanceAddress := getBalanceCmd.String("address", "", "create block and package dtata")
	switch os.Args[1] {

	case "create":
		err := crateBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "get":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printf":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if crateBlockchainCmd.Parsed() {
		if *genenisAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Println("创建创世区块并且存储到数据库")
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
	if getBalanceCmd.Parsed() {
		if *balanceAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Printf("查询%s 的余额",balanceAddress)
	}

}
