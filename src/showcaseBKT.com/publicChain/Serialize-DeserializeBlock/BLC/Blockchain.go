package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
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

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

//迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

// 遍历输出所有区块的信息
func (blockchian *Blockchain) Printchain() {
	fmt.Println("遍历输出所有区块的信息....")
	blockchainIterator := blockchian.Iterator()

	for {
		block := blockchainIterator.Next()
		//fmt.Println(block)
		//fmt.Println("============START==============================")
		fmt.Printf("Hash :%x \n", block.Hash)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Height:%d \n", block.Height)
		//fmt.Printf("Timestamp		:%s \n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
		//fmt.Printf("Nonce			:%d \n", block.Nonce)

		for _, transcation := range block.Txs {
			//fmt.Println("**************************")
			fmt.Printf("####交易id:		%x\n", transcation.TxHash)
			//fmt.Println("\t-------Vins:")
			for _, in := range transcation.Vins {
				fmt.Printf("\tvin txid        :%x\n", in.TxHash)
				fmt.Printf("\tvin voutIndex   :%d\n", in.VoutIndex)
				fmt.Printf("\tvin 付款人的公钥:%x\n", in.PublicKey)
				fmt.Printf("\tvin 签名:%x\n", in.Signature)
			}

		//	fmt.Println("\t--------Vouts:")
			for _, out := range transcation.Vouts {
				fmt.Printf("\tvout 收款人的公钥:%x\n",( out.Ripemd160Hash))
				fmt.Printf("\tvout 收款的金额数量      :%d\n", out.Value)
			}
		//	fmt.Println("**************************")
		}
		fmt.Println("")
		var hashBigInt big.Int
		//是否到达创世区块
		hashBigInt.SetBytes(block.PrevBlockHash)
		if (hashBigInt.Cmp(big.NewInt(0)) == 0) {
			break;
		}

	}

}

//从一个地址对应的TXoutput未花费
// 多笔转帐时要把准备打包的txs也要传进来一同计算
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transcation) []*UTXO {

	var unUTXOs []*UTXO
	//某个 tx 对应它已消费的input的VoutIndex, 001->int{0,1}
	spentTXOutputs := make(map[string][]int)
	fmt.Printf("开始计算%s 地址对应的未花费的TXO\n",address)
	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//公钥hash
				publicKey := Base58Decode([]byte(address))
				ripemd160hash := publicKey[1 : len(publicKey)-4]

				if in.UnlockRipedm160Hash(ripemd160hash) {
					fmt.Println("打印所有txs ,tx hash- in hash-inindex:", hex.EncodeToString(tx.TxHash), ",in hash", hex.EncodeToString(in.TxHash), in.VoutIndex);
				}
			}
		}
	}

	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				//公钥hash
				publicKey := Base58Decode([]byte(address))
				ripemd160hash := publicKey[1 : len(publicKey)-4]

				if in.UnlockRipedm160Hash(ripemd160hash) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.VoutIndex)
				}
			}
		}
	}
	//fmt.Println("=======每个Txs中 已消费的 output===============")
	for key, indexArray := range spentTXOutputs {
		for _, voutIndex := range indexArray {
			fmt.Println("已消费的output key - index:", key, voutIndex);
		}
	}

	for _, tx := range txs {
	Work1:
		for index, out := range tx.Vouts {
			//是否已经被花费
			if out.UnLockScriptPubKeyWithAddress(address) {
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					fmt.Println("inser进UTXOs  len(spentTXOutputs) == 0 :", hex.EncodeToString(tx.TxHash), index, out);
					unUTXOs = append(unUTXOs, utxo)
				} else {
					var isSpentUTXO bool
					for txHash, indexArray := range spentTXOutputs {
						for _, voutIndex := range indexArray {
							if index == voutIndex && txHash == hex.EncodeToString(tx.TxHash) {
								isSpentUTXO = true
								continue Work1
							}
						}
					} // end spentTXOutputs
					if isSpentUTXO == false {
						fmt.Println("inser进UTXOs  ok", hex.EncodeToString(tx.TxHash), index, out);
						utxo := &UTXO{tx.TxHash, index, out}
						unUTXOs = append(unUTXOs, utxo)
					}
				}
			} //end if
		}
	}//end range txs


	blockIterator := blockchain.Iterator()
	for {
		block := blockIterator.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]

			if tx.IsCoinbaseTransaction() == false {

				for _, in := range tx.Vins {
					//是否能够解锁 	//公钥hash
					publicKey := Base58Decode([]byte(address))
					ripemd160hash := publicKey[1 : len(publicKey)-4]

					if in.UnlockRipedm160Hash(ripemd160hash) {
						//将transcation id(byte array) 转成string
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.VoutIndex)

					}
				}
			} //end if

			// Vouts
		work:
			for index, out := range tx.Vouts {
				//是否已经被花费

				if out.UnLockScriptPubKeyWithAddress(address) {

					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							//key->[]index
							//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]
							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {
								for _, voutIndex := range indexArray {
									if index == voutIndex && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)

						}
					}
				}else{

				}
			} //end Vouts for
		} //end tx for

		var hashBigInt big.Int
		//是否到达创世区块
		hashBigInt.SetBytes(block.PrevBlockHash)

		if (hashBigInt.Cmp(big.NewInt(0)) == 0) {
			break;
		}
	} // end for

	for _, outx := range unUTXOs {
		fmt.Printf ("打印未花费的UTXO index:%d ,TxHash %x\n:", outx.Index, (outx.TxHash))
	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int, txs []*Transcation) (int64, map[string][]int) {

	//查看未花费
	utxos := blockchain.UnUTXOs(from, txs)

	//交易id对应未消费的txoutput的index
	spendableUTXO := make(map[string][]int)

	//2. 遍历utxos
	var value int64 //统计【unspentOutputs】对应的txoutput未花费总量

	for _, utxo := range utxos {
		value = value + utxo.OutPut.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= int64(amount) {
			break
		}
	} //end for

	if value < int64(amount) {
		fmt.Printf("%s 's fund is 不足\n", from)
		os.Exit(1)
	}
	return value, spendableUTXO
}

// 挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {

	//main send -from [\"helloggggggggg\",\"apple\",\"java\",\"golang\"] -to [\"bbc\",\"btc\",\"bkc\",\"blc\"] -amount [\"10\",\"20\"]

	//1.建立一笔交易
	//1. 通过相关算法建立Transaction数组
	var txs []*Transcation
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
	}

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			hash := b.Get([]byte(TIP))

			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	//2.创建区块
	newBlock := NewBlock(txs, block.Height+1, block.Hash)

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

// 返回Blockchain对象
func BlockchainObject() *Blockchain {
	if DBExists() == false {
		fmt.Println("数据不存在..")
		//log.Panic("数据不存在..")
		os.Exit(1)
	}

	//1.尝试打开或是创建数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//2.读取最未尾的区块
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(TIP))
		return nil
	})

	return &Blockchain{tip, db}
}

//1.创建一个带有创或区块的区块键
func CreateBlockchainWithGenesisBlock(genesisAddress string) *Blockchain {
	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}


	//return &Blockchain{[]*Block{NewGenesisBlock()}}
	//1.尝试打开或是创建数据库
	//如果数据存在 就打开，否则创建一个数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tip [] byte //获取最后一个区块hash
	//2 db.update更新数据
	//2.1 表是否存在，如果不存在需要创建表
	err = db.Update(func(tx *bolt.Tx) error {
		//判断一张表是否存在于数据库中
		b := tx.Bucket([]byte(blocksBucket))

		//表不存
		if b == nil {
			//fmt.Println(" No existing blockchain found. create a new one.")
			fmt.Println("正在创建创世区块.......")
			//表不存认为区块为空，需要首次创建创世区块

			//创建创世区块的交易对象
			cbtx := NewCoinbaseTransaction(genesisAddress)
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

// 查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {
	utxos := blockchain.UnUTXOs(address, []*Transcation{})
	var amount int64
	for _, utxo := range utxos {
		//fmt.Println(utxo.TxHash, utxo.Index, utxo.OutPut.Value)
		amount = amount + utxo.OutPut.Value
	}
	return amount
}

func (blockchain *Blockchain) SignTranscation(tx *Transcation, privKey ecdsa.PrivateKey) {

	if tx.IsCoinbaseTransaction() {
		return
	}

	//vin 所对应的 交易集
	prevTXs := make(map[string]Transcation)
	for _, vin := range tx.Vins {
		prevTX, err := blockchain.FindTransction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}
//通过 ID 找到一笔交易（这需要在区块链上迭代所有区块）
func (bc *Blockchain) FindTransction(ID []byte) (Transcation, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return Transcation{}, nil
}
