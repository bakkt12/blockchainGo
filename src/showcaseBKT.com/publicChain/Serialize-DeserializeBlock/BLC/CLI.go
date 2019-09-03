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

//直接打印usage信息
func (cli *CLI) printUsage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\taddresslist		\t--输出所有钱包地址")
	fmt.Println("\tcreatewallet		\t--创建钱包")
	fmt.Println("\tcreateblockchain \t-address \"创建创世区块的地址...\"")
	fmt.Println("\tgetbalance			\t-address  \"要查询某一个账号的余额......\".")
	fmt.Println("\tprintf				\t 输出所有区块的数据.........")
	fmt.Println("\tsend				\t-from \"转账源地址...\" -to \"转账目的地地址...\"  -amount \"转账金额......\"")
}

func (cli *CLI) Run() {
	//判断终端参数的个数 如果没有参数直接打印usage信息
	cli.validateArgs()

	createwalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	addresslistCmd := flag.NewFlagSet("addresslist", flag.ExitOnError)

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
	case "createwallet":
		err := createwalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addresslist":
		err := addresslistCmd.Parse(os.Args[2:])
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
			//cli.printUsage()
			os.Exit(1)
		}
		fromAddress := JSONtoArray(*sendFrom)
		toAddress := JSONtoArray(*sendTo)
		amount := JSONtoArray(*sendAmount)

		for index, fromAd := range fromAddress {
			if IsValidForAddress([]byte(fromAd)) == false || IsValidForAddress([]byte(toAddress[index])) == false {
				fmt.Printf("地址无效 %s ,%s \n", fromAd, toAddress[index])
				//cli.printUsage()
				os.Exit(1)
			}
		}

		fmt.Printf("from %s\n", fromAddress)
		fmt.Printf("to: %s\n", toAddress)
		fmt.Printf("amount: %s\n", amount)
		cli.Send(fromAddress, toAddress, amount)
	}
	if printChainCmd.Parsed() {
		fmt.Println("开始输出所有区块的数据........")
		cli.PrintChain()
	}

	if crateBlockchainCmd.Parsed() {
		fmt.Println("开始创建创世区块....")
		if IsValidForAddress([]byte(*flagCreateBlockchainWithAddress )) == false {
			fmt.Println("创建创世区块地址无效....")
			//	cli.printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}

	if getbalanceCmd.Parsed() {
		fmt.Printf("[CLI] 开始查询 %s 地址余额........\n", *getbalanceWithAdress)
		if IsValidForAddress([]byte(*getbalanceWithAdress )) == false {
			fmt.Println("[CLI] 查询地址无效....")
			//cli.printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceWithAdress)
	}

	if createwalletCmd.Parsed() {
		fmt.Println("开始创建创钱包....")
		cli.createWallet()
	}

	if addresslistCmd.Parsed() {
		fmt.Println("开始输出创钱包地址....")
		cli.getaddresslists()
	}

}
