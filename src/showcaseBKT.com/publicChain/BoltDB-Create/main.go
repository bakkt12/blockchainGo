package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func main() {
	// -------------数据库创-----------------------
	//如果数据存在 就打开，否则创建一个数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close();
	err = db.Update(func(tx *bolt.Tx) error {

		//判断一张表是否存在于数据库中
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			fmt.Println(" No existing blockchain found. create a new one.")

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			// 存储数据
			err = b.Put([]byte("yehaoning"), []byte("http://bakkt.org"))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("yjc"), []byte("http://yjc.org"))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

}
