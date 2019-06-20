package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	//BC *Blockchain
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
	existDB := DBExists()
	if existDB == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	blockchian := BlockchainObject();
	defer blockchian.DB.Close()

	blockchian.Printchain()
}

//直接打印usage信息
func (cli *CLI) printUsage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain  \t-address \"创建创世区块的地址...\"")
	fmt.Println("\tgetbalance			\t-address  \"要查询某一个账号的余额......\".")
	fmt.Println("\tprintf				\t-print 输出所有区块的数据.........")
	fmt.Println("\tsend				\t-from \"转账源地址...\" -to \"转账目的地地址...\"  -amount \"转账金额......\"")
}

func (cli *CLI) Run() {
	//判断终端参数的个数 如果没有参数直接打印usage信息
	cli.validateArgs()
	//创建区块
	crateBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	flagCreateBlockchainWithAddress := crateBlockchainCmd.String("address", "", "创建创世区块的地址...")
	//打印
	printChainCmd := flag.NewFlagSet("printf", flag.ExitOnError)
	//查询
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	//转帐
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendBlockCmd.String("from", "", "转账源地址...")
	sendTo := sendBlockCmd.String("to", "", "转账目的地地址...")
	sendAmount := sendBlockCmd.String("amount", "", "转账金额......")

	switch os.Args[1] {

	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := crateBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
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

	if sendBlockCmd.Parsed() {
		fmt.Printf("开始转帐， 转账源地址%s -> 转账目的地址:%s,转账金额:%s \n", *sendFrom, *sendTo, *sendAmount)
		if *sendFrom == "" || *sendTo == "" || *sendAmount == "" {
			fmt.Println("转账源地址,转账目的地址, 转账金额 不能为空....")
			cli.printUsage()
			os.Exit(1)
		}
		fromAddress := JSONtoArray(*sendFrom)
		toAddress := JSONtoArray(*sendTo)
		amount := JSONtoArray(*sendAmount)

		fmt.Printf("from %s\n", fromAddress)
		fmt.Printf("to: %s\n", toAddress)
		fmt.Printf("amount: %s\n", amount)
		cli.send(fromAddress, toAddress, amount)
	}
	if printChainCmd.Parsed() {
		fmt.Println("开始输出所有区块的数据........")
		cli.PrintChain()
	}

	if crateBlockchainCmd.Parsed() {
		fmt.Println("开始创建创世区块....")
		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("地址不能为空....")
			cli.printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}

	if getbalanceCmd.Parsed() {
		fmt.Printf("开始查询%s地址余额........\n",*getbalanceWithAdress)
		if *getbalanceWithAdress == "" {
			fmt.Println("查询地址不能为空....")
			cli.printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceWithAdress)
	}
}

//创建创世区块
func (cli *CLI) createGenesisBlockchain(genesis string) {
	if DBExists() {
		fmt.Println("创世区块已经存在")
		os.Exit(1)
	}
	CreateBlockchainWithGenesisBlock(genesis)
}

//转帐
func (cli *CLI) send(from []string, to []string, amount []string) {
	//if DBExists() == false {
	//	fmt.Println("数据不存在...")
	//	os.Exit(1)
	//}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amount)
}
// 查询余额
func (cli *CLI) getBalance(address string) {
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount :=blockchain.GetBalance(address)
	fmt.Printf("%s 一共有%d个Token\n",address,amount)
}
