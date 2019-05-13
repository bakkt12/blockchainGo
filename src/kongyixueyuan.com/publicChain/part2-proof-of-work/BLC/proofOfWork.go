package BLC

import (
	"fmt"
	"github.com/imroc/biu"
	"math/big"
)

const targetBits = 34

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func (pow *ProofOfWork) Run() {
	fmt.Printf("Mining the block containing \"%s\" \n ", pow.block.Data)
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	fmt.Println(target)
	target.Lsh(target, uint(256-targetBits))
	fmt.Printf("%d \n", target)

	pow := &ProofOfWork{block, target}
	return pow
}
