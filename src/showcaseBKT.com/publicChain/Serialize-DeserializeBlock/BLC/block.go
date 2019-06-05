package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	//时间戳 创建区块时的时间
	Timestamp int64
	//上一个区块Hash
	PrevBlockHash []byte
	//交易数据
	//旧逻辑 交易数据用Data []byte ->Transcation []Transcation
	Transcation []*Transcation
	//当前区块Hash
	Hash []byte
	//Nonce 随机数
	Nonce int

	//nBits int                // 块目标值
}

//Not used!
/*func (block *Block) SetHash() {
	// 1. 时间转字节数组
	//(1) init64转化为字符串  int64到string
	timeString := strconv.FormatInt(block.Timestamp, 2)
	//（2）字符串转化为字节数组
	timestamp := []byte(timeString)
	//  2.以字节数组方式拼接 上一区块hash,当前区块数据，时间戳
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	//3.数据进行256hash
	hash := sha256.Sum256(headers)
	//4.hash设置为当前 hash.
	block.Hash = hash[:]
}*/

/**
  产生新的区块工厂方法，
*/
func
NewBlock(transcation []*Transcation, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, transcation, []byte{}, 0}
	//将block作为参数 创建一个pow对象
	pow := NewProofOfWork(block)

	//执行一次工作量证明
	noce, hash := pow.Run()
	//设置区块Hash
	block.Hash = hash[:]
	//设置Nonce
	block.Nonce = noce

	isValid := pow.validate();
	fmt.Printf("newBlock %c \n\n", isValid)

	return block
}

//将block 对象序列化为[]byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer;
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//将区块里面所有的交易id拼接生成Hash
func (b *Block) HashTranscation() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transcation {
		txHashes = append(txHashes, tx.ID)
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

//反序列化
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewBuffer(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//创建创世区块
func NewGenesisBlock(conbase *Transcation) *Block {
	return NewBlock([]*Transcation{conbase}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
