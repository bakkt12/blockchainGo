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

	//Read-only transactions
	//查询数据
	err = db.View(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blocksBucket))
		valueByte := b.Get([]byte("yjc"))

		fmt.Printf("%s \n", valueByte)
		valueByte = b.Get([]byte("yehaoning"))
		fmt.Printf("%s \n", valueByte)
		return nil
	})

}
