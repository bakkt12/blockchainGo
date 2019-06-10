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

//判断终端参数的个
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) PrintChain() {

	//检查是否有数据存在
	existDB := dbExists()
	fmt.Printf(" dbExists() %T:\n", existDB)
	fmt.Printf("%d\n", existDB) // {1 2}

	fmt.Printf("%+v\n", existDB) // {x:1 y:2}

	fmt.Printf("%#v\n", existDB) // ma
	if existDB == false {
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
	tx1 := NewUTXOTransaction("yhn", "BAKKT", 5, cli.BC, noPackageTxs)
	noPackageTxs = append(noPackageTxs, tx1)

	tx2 := NewUTXOTransaction("yhn", "YE", 5, cli.BC, noPackageTxs)
	noPackageTxs = append(noPackageTxs, tx2)
	////
	tx3 := NewUTXOTransaction("yhn", "LY", 5, cli.BC, noPackageTxs)
	noPackageTxs = append(noPackageTxs, tx3)

	cli.BC.MineBlock([]*Transcation{tx1, tx2, tx3})
}

func (cli *CLI) addBlock(data string) {
	cli.SendToken();
}

//直接打印usage信息
func (cli *CLI) printUsage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain  \t-address \"address\"")
	fmt.Println("\tgetbalance			\t-address  \"address\".")
	fmt.Println("\tprintf				\t-print all the blocks of the blockchain.")
	fmt.Println("\tsend				\t-from -to -amount")
}

func (cli *CLI) Run() {
	//判断终端参数的个数 如果没有参数直接打印usage信息
	cli.validateArgs()
	//创建区块
	crateBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	genenisAddress := crateBlockchainCmd.String("address", "", "create block and package dtata")
	//打印
	printChainCmd := flag.NewFlagSet("printf", flag.ExitOnError)
	//查询
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	balanceAddress := getBalanceCmd.String("address", "", "create block and package dtata")

	//转帐
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "源地址")
	sendTo := sendCmd.String("to", "", "目标地址")
	sendAmount := sendCmd.String("amount", "", "转帐的额度")

	switch os.Args[1] {

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := crateBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
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
		cli.addBlock("")
	}

	if sendCmd.Parsed() {
		fmt.Printf("from:%s, to:%s,amount %s \n", *sendFrom, *sendTo, *sendAmount)
		if *sendFrom == "" || *sendTo == "" || *sendAmount == "" {
			fmt.Println("null---")
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Println("json ->array[]")
		fromAddress:=JSONtoArray(*sendFrom)
		toAddress:=JSONtoArray(*sendTo)
		sendAmount:=JSONtoArray(*sendAmount)

		fmt.Printf("from %s\n",fromAddress)
		fmt.Println("to: %s\n",toAddress)
		fmt.Println("amount %s\n",sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}

	if getBalanceCmd.Parsed() {
		if *balanceAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Printf("查询%s 的余额", balanceAddress)
	}

}
