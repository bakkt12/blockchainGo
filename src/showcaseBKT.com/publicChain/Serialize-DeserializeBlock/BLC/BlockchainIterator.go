package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte   //当前正在遍历的区块hash
	DB          *bolt.DB //数据库
}


//获取前一个迭代器
func (blockchainIterator *BlockchainIterator) Next() *Block {
	//var nextHash []byte
var block *Block

	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)

			//  获取到当前迭代器里面的currentHash所对应的区块
			block= DeserializeBlock(currentBlockBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.CurrentHash = block.PrevBlockHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	//return &BlockchainIterator{nextHash, blockchainIterator.DB}
	return block;
}

