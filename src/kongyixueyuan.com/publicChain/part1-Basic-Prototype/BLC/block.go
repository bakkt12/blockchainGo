package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	//时间戳 创建区块时的时间
	Timestamp int64
	//上一个区块Hash
	PrevBlockHash []byte
	//交易数据
	Data [] byte
	//当前区块Hash
	Hash [] byte
}

func (block *Block) SetHash() {
	// 1. 时间转字节数组
	//(1) init64转化为字符串
	timeString := strconv.FormatInt(block.Timestamp, 2)
	//（2）字符串转化为字节数组
	timestamp := []byte(timeString)
    //  2.以字节数组方式拼接 上一区块hash,当前区块数据，时间戳
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	//3.数据进行256hash
	hash := sha256.Sum256(headers)
	//4.hash设置为当前 hash.
	block.Hash = hash[:]
}

/**
  产生新的区块工厂方法，
 */
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), PrevBlockHash: prevBlockHash, Data: []byte(data), Hash: []byte{}}
	//fmt.Println(block)
	block.SetHash()
	return block;
}

func NewGenesisBlock() *Block {
	return NewBlock("Genenis Block", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
