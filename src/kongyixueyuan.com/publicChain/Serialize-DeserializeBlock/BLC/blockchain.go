package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//数据库名
const dbFile = "blockchain.db"

//表
const blocksBucket = "blocks"

type Blockchain struct {
	//Blocks []*Block //存储有序的区块
	Tip []byte   //区块键里最后一个区块的Hash
	DB  *bolt.DB //数据库
}

//新增区块
func (blockchain *Blockchain) AddBlock(data string) {
	////1.创新区块
	//preBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	//newBlock := NewBlock(data, preBlock.Hash)
	//
	////2 加到blockchain里面
	//blockchain.Blocks = append(blockchain.Blocks, newBlock)

	//=====================================================
	//1.创建区块
	newBlock := NewBlock(data, blockchain.Tip)

	//2. update 数据库
	err := blockchain.DB.Update(func(tx *bolt.Tx) error {
		//2.1获取表
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		//更新 l对应的Hash
		err = b.Put([]byte("L"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		//将最新的区块存储到blockchain tip中。
		blockchain.Tip = newBlock.Hash;
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

//创建一个带有创或区块的区块键
func NewBlockchain() *Blockchain {

	var tip [] byte //获取最后一个区块hash
	//return &Blockchain{[]*Block{NewGenesisBlock()}}
	//1.尝试打开或是创建数据库
	//如果数据存在 就打开，否则创建一个数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	//2 db.update更新数据
	//2.1 表是否存在，如果不存在需要创建表
	err = db.Update(func(tx *bolt.Tx) error {
		//判断一张表是否存在于数据库中
		b := tx.Bucket([]byte(blocksBucket))

		//表不存
		if b == nil {
			fmt.Println(" No existing blockchain found. create a new one.")
			//创建创世区块
			genesisBlock := NewGenesisBlock()
			//创建表
			b, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			//创世区块序列之后的数据存储到表中
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("L"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = genesisBlock.Hash
		} else { //表存在
			//KEY:L ,value为最后一个区块hash
			tip = b.Get([]byte("L"))
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	//2.2 创建创世区块，序列化
	//2.3 创世区块的hash作为key block序列化数据做为value 存在表中
	//2.4 设置一个key L,将hash作为value存进数据库中。作为最后一个 区块

	return &Blockchain{tip, db}
}
