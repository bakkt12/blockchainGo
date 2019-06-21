package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsildy = 50

/**
交易 - 幕后细节
{
  "version": 1,
  "locktime": 0,
  "vin": [
    {
      "txid":"7957a35fe64f80d234d76d83a2a8f1a0d8149a41d81de548f0a65a8a999f6f18",
      "vout": 0,
      "scriptSig": "3045022100884d142d86652a3f47ba4746ec719bbfbd040a570b1deccbb6498c75c4ae24cb02204b9f039ff08df09cbe9f6addac960298cad530a863ea8f53982c09db8f6e3813[ALL] 0484ecc0d46f1918b30928fa0e4ed99f16a0fb4fde0735e7ade8416ab9fe423cc5412336376789d172787ec3457eee41c04f4938de5cc17b4a10fa336a8d752adf",
      "sequence": 4294967295
    }
 ],
  "vout": [
    {
      "value": 0.01500000,
      "scriptPubKey": "OP_DUP OP_HASH160 ab68025513c3dbd2f7b92a94e0581f5d50f654e7 OP_EQUALVERIFY OP_CHECKSIG"
    },
    {
      "value": 0.08450000,
      "scriptPubKey": "OP_DUP OP_HASH160 7f9b1a7fb68d60c536c2fd8aeaa53a8f3cc025a8 OP_EQUALVERIFY OP_CHECKSIG",
    }
  ]
}


 */
type Transcation struct {
	//Vesion string 未用到
	TxHash []byte      // 交易id
	Vins   []*TXInput  // 交易输入
	Vouts  []*TXOutput //交易输出
	//	Lock_time int64 未用到
}

// 1.判断当前交易是否是 coinbase tx
func (tx *Transcation) IsCoinbaseTransaction() bool {
	return tx.Vins[0].VoutIndex == -1 && len(tx.Vins[0].TxHash) == 0
}



//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func NewCoinbaseTransaction(address string) *Transcation {

	fmt.Sprintf("NewCoinbaseTransaction to '%s'", address)

	//创建创世的输入
	txin := &TXInput{[]byte{}, -1, "Genesis Data"}
	//创建输出
	txout := &TXOutput{subsildy, address}
	//创建交易
	tx := Transcation{nil, []*TXInput{txin}, []*TXOutput{txout}}
	tx.HashTransaction()
	return &tx
}

//建立交易
//2. 转账时产生的Transaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain) *Transcation {

	fmt.Printf("Create a NewSimpleTransaction..from:%s->to:%s \n", from, to)

	//1.找到有效的可用的交易输出数据模型
	//查询出未花费的输出  (int,map[string][]int)

	money, spendableUTXODic := blockchain.FindSpendableUTXOS(from, amount)

	fmt.Println("NewSimpleTransaction")
	fmt.Println(spendableUTXODic)
	fmt.Println(amount)

	//输入
	var txIntputs []*TXInput
	//输出
	var txOutputs []*TXOutput

	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txIntput := &TXInput{txHashBytes, index, from}
			txIntputs = append(txIntputs, txIntput)
		}
	}

	//建立输出转帐
	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)

	//建立输出，找零
	txOutput = &TXOutput{int64(money) - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)

	//创建交易
	tx := &Transcation{[]byte{}, txIntputs, txOutputs}

	//设置hash值
	tx.HashTransaction()
	return tx
}

//设置交易 hash，将 tx交易序列化之后字节数组生成256Hash
func (tx *Transcation) HashTransaction() {
	var encoded bytes.Buffer;
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(encoded.Bytes())
	tx.TxHash = hash[:]
}
