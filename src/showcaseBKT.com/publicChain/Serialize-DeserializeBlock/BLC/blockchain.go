package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"time"
)

//数据库名
const dbFile = "blockchain.db"

//表
const blocksBucket = "blocks"

//创世区块的数据信息
const genesisCoinbaseData = "The times 03/Jan/2009 Chancellor on brink of second bailout for banks"

const TIP = "L"

type Blockchain struct {
	//Blocks []*Block //存储有序的区块
	Tip []byte   //区块键里最后一个区块的Hash
	DB  *bolt.DB //数据库
}

/**
找到包含当前用户未花费的输出的所有交易集合
返回交易数组集
 */
func (bc *Blockchain) FindUnspentTranscations(address string) [] Transcation {

	//存储未花费输出的交易
	var unspentTXs []Transcation

	//key:transcation id- output index
	//存储交易所对应已花费的(vinput里的)对应的output的index
	spentTXOs := make(map[string][]int)

	blockchainIterator := bc.Iterator()

	var hashBigInt big.Int
	fmt.Println("========FindUnspentTranscations=============")
	for {
		err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			//通过hash获取到区块字节数组
			currentBlockBytes := b.Get([]byte(blockchainIterator.CurrentHash))
			currentBlock := DeserializeBlock(currentBlockBytes)

			fmt.Printf("PrevBlockHash:%x \n", currentBlock.PrevBlockHash)
			fmt.Printf("Hash:%x \n", currentBlock.Hash)
			fmt.Printf("Nonce:%d \n", currentBlock.Nonce)
			fmt.Printf("Timestamp:%s \n\n", time.Unix(currentBlock.Timestamp, 0).Format("2006-01-02 15:04:05"))
			for _, transcation := range currentBlock.Transcation {
				fmt.Printf("TranscationHash %x \n", transcation.ID)

				//将transcation id(byte array) 转成string
				txId := hex.EncodeToString(transcation.ID)

			Outputs:
				for outIdx, out := range transcation.Vout {
					//是否已经被花费
					/**
					  比特币应用可以使用一些策略来满足付款需要：组合若干小的个体，算出准确的找零；或者使用一个比交易值大的个体然后进行找零。
					   vout是否在区块的 vint中
					 */
					if spentTXOs[txId] != nil {
						for _, spentOut := range spentTXOs[txId] {
							if spentOut == outIdx {
								//终止当前for循环 ，从Outputs标签处for继续执行
								//GO的 continue只是停止执行continue之后的语句，但会继续在当前for中执行
								continue Outputs
							}
						} //end for
					}

					if out.CanBeUnlockedWith(address) {
						unspentTXs = append(unspentTXs, *transcation)
					}
				} //end for Vout

				if transcation.IsCoinbase() == false {
					for _, in := range transcation.Vin {
						if in.CanUnlockOutputWith(address) {
							inTxId := hex.EncodeToString(in.Txid)
							spentTXOs[inTxId] = append(spentTXOs[inTxId], in.VoutIndex)
						}
					}
				}
			} //end for transcation
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		//获取下一个迭代器
		blockchainIterator = blockchainIterator.Next()

		//是否到达创世区块
		hashBigInt.SetBytes(blockchainIterator.CurrentHash)
		if (hashBigInt.Cmp(big.NewInt(0)) == 0) {
			break;
		}
	} //end for
	return unspentTXs
}

//查找可用的未消费的输出信息
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	//{"1111":[1,2,3]}
	//交易id对应未消费的txoutput的index
	unspentOutputs := make(map[string][]int)
	//查看未花费
	unspentTXs := bc.FindUnspentTranscations(address)
	accumulated := 0 //统计【unspentOutputs】对应的txoutput未花费总量

Work:
	for _, tx := range unspentTXs {
		txId := hex.EncodeToString(tx.ID)
		for index, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.CAmount
				unspentOutputs[txId] = append(unspentOutputs[txId], index)
				if accumulated >= amount {
					break Work
				}
			}
		} //end for
	} //end for

	return accumulated, unspentOutputs
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
	/*newBlock := NewBlock(data, blockchain.Tip)

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
	}*/
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

			//创建创世区块的交易对象
			cbtx := NewCoinbaseTx("yhn", genesisCoinbaseData)
			genesisBlock := NewGenesisBlock(cbtx)
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
