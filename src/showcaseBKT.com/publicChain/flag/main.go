package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func testOs() {
	args := os.Args;
	fmt.Println(args)

	wordPtr := flag.String("word", "food", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("bool", false, "a bool")

	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("tail:", flag.Args())
}

func testNewFlagSet() {

	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Add")

	//commmand:main addBlock -data btc
	fmt.Println(os.Args[0])  // main
	fmt.Println(os.Args[1])  // addblock
	fmt.Println(os.Args[2:]) //-data btc
	switch os.Args[1] {
	case "addBlock":
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
		fmt.Println("No addblock and printchain")
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}

		fmt.Println("Data:" + *addBlockData)
	}

	if printChainCmd.Parsed() {
		fmt.Println("printchaint!")
	}

}
func main() {

}
