package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
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
func (cli *CLI) Run() {
	//判断终端参数的个数 如果没有参数直接打印usage信息
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block")
	fmt.Printf("os args %s \n", os.Args)

	if len(os.Args) < 2 {

	}
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
		//fmt.Println("No addblock and printchain!")
		flag.Usage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		fmt.Printf("Data:" + *addBlockData)
	}

}
