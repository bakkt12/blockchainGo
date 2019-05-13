package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 24

var (
	//定义最大
	maxNonce = math.MaxInt64
)

type ProofOfWork struct {
	block  *Block
	target *big.Int //区块难度值
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		}, []byte{},
	)

	return data;
}

func (pow *ProofOfWork) Run() (int, []byte) {
	fmt.Printf("Mining the block containing \"%s\" \n ", pow.block.Data)
	var hashInt big.Int
	var hash [32]byte

	nonce := 0;
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r %x", hash)
		hashInt.SetBytes(hash[:])
		//  hasint < pow .target -1
		//  hasint = pow .target -1  0
		//  hasint > pow .target -1   1
		if hashInt.Cmp(pow.target) == -1 {
			break
		}else{
			nonce++
		}
	}
	return nonce, []byte{};
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	/**
	左移多少位，相当于（8-5）2的3次方
	00000001  =1
	00100000  //移5位
	 */
	fmt.Printf("sart %b \n", target)
	target.Lsh(target, uint(256-targetBits))
	fmt.Printf("after %b \n", target)

	pow := &ProofOfWork{block, target}
	return pow
}
