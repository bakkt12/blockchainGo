package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	//交易数据
	Data []byte
	//当前区块Hash
	Nonce int
}

//Block 对象序化成[]byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//字节数组反序列化成Block
func  DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewBuffer(d))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return &block
}
