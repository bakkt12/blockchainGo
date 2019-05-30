package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte   //当前正在遍历的区块hash
	DB          *bolt.DB //数据库
}

//迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

//获取前一个迭代器
func (bi *BlockchainIterator) Next() *BlockchainIterator {
	var nextHash []byte

	err := bi.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		currentBlockBytes := b.Get([]byte(bi.CurrentHash))

		currentBlock := DeserializeBlock(currentBlockBytes)

		nextHash = currentBlock.PrevBlockHash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &BlockchainIterator{nextHash, bi.DB}
}
