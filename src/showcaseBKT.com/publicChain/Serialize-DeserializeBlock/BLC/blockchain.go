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

const TIP = "L"

type Blockchain struct {
	//Blocks []*Block //存储有序的区块
	Tip []byte   //区块键里最后一个区块的Hash
	DB  *bolt.DB //数据库
}

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
		//fmt.Printf("next hash %x:\n" ,nextHash)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &BlockchainIterator{nextHash, bi.DB}
}

//新增区块
func (blockchain *Blockchain) AddBlock(data string) {
	/*##################旧逻辑#############################################**/
	//1.创新区块
	//preBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	//newBlock := NewBlock(data, preBlock.Hash)
	//2 加到blockchain里面
	//blockchain.Blocks = append(blockchain.Blocks, newBlock)
	/*##################################################################**/

	//用数据库来存放整个区块
	//1.创建区块
	newBlock := NewBlock(data, blockchain.Tip)

	//2. update 数据库
	err := blockchain.DB.Update(func(tx *bolt.Tx) error {

		//2.1获取表
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		//更新 l对应的Hash
		err = b.Put([]byte(TIP), newBlock.Hash)
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

//创取当前最新的区块
func LastBlockchain() *Blockchain {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(TIP))
		return nil
	})
	return &Blockchain{tip, db}
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
	//defer db.Close()

	//2 db.update更新数据
	//2.1 表是否存在，如果不存在需要创建表
	err = db.Update(func(tx *bolt.Tx) error {
		//判断一张表是否存在于数据库中
		b := tx.Bucket([]byte(blocksBucket))

		//表不存
		if b == nil {
			fmt.Println(" No existing blockchain found. create a new one.")
			//表不存认为区块为空，需要首次创建创世区块
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

			err = b.Put([]byte(TIP), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = genesisBlock.Hash
		} else {
			//表存在
			//KEY:L ,value为最后一个区块hash
			tip = b.Get([]byte(TIP))
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
