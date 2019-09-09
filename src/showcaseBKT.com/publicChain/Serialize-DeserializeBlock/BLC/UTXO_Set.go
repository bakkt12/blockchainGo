package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

const utxoTableName = "utxoTableName"

type UTXOSet struct {
	Blockchain *Blockchain
}

//重置数据库的表
func (utxoSet *UTXOSet) ResetUTXOSet() {

	//1. update 数据库
	err := utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		//1.2获取表
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			tx.DeleteBucket([]byte(utxoTableName))
			b, _ = tx.CreateBucket([]byte(utxoTableName))
			if b != nil {
				txOutPutMap :=utxoSet.Blockchain.FindUTXOMap()
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
