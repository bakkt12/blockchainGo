package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsildy = 12.5

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

//1. Transaction 创建分两种情况
//1. 创世区块创建时的Transaction
func NewCoinbaseTransaction(address string) *Transcation {

	fmt.Sprintf("NewCoinbaseTransaction to '%s'", address)

	//创建创世的输入
	txin := &TXInput{[]byte{}, -1, nil, []byte{}}
	//创建输出
	txout := NewTXOutput(10, address)
	//创建交易
	tx := Transcation{[]byte{}, []*TXInput{txin}, []*TXOutput{txout}}
	tx.HashTransaction()
	return &tx
}

//建立交易
//2. 转账时产生的Transaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transcation) *Transcation {
	fmt.Printf("Create a new transaction from:%s -> to:%s, amount:%d \n", from, to, amount)
	//1.找到有效的可用的交易输出数据模型
	//查询出未花费的输出  (int,map[string][]int)

	//钱包获取 原生公钥
	wallets, _ := NewWallets()
	wallet := wallets.GetWallet(from)

	money, spendableUTXODic := blockchain.FindSpendableUTXOS(from, amount, txs)
	fmt.Printf("[Transcation.go]>from:%s 找到可用 monety:%d \n", from, money)

	//输入
	var txIntputs []*TXInput

	//输出
	var txOutputs []*TXOutput

	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txIntput := &TXInput{txHashBytes, index, nil, wallet.PublicKey}
			txIntputs = append(txIntputs, txIntput)
		}
	}

	//建立输出转帐
	txOutput := NewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	//建立输出，找零
	txOutput = NewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	//创建交易
	tx := &Transcation{[]byte{}, txIntputs, txOutputs}

	//设置hash值
	tx.HashTransaction()

	//进行签名
	blockchain.SignTranscation(tx, wallet.PrivateKey)
	return tx
}

//方法接受一个私钥和一个之前交易的 map
func (tx *Transcation) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transcation) {
	if (tx.IsCoinbaseTransaction()) {
		return
	}

	txCopy := tx.TrimmedCopy()
	//每一个输入都是分开签名，交易中可以是包含不同地址的输入

	//接下来，我们检查每个输入中的签名
	for inID, vin := range txCopy.Vins {
		//找到input对应的是上一个transcation,
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		//PubKey 被设置为所引用输出的 PubKeyHash
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.VoutIndex].Ripemd160Hash
		// 方法对交易进行序列化，并使用 SHA-256 算法进行哈希
		//哈希后的结果就是我们要签名的数据
		txCopy.TxHash = txCopy.Hash()
		//重置为null不影响后面的迭代
		txCopy.Vins[inID].PublicKey = nil

		//通过 privKey 对 txCopy.ID 进行签名。一个 ECDSA 签名就是一对数字，
		// 我们对这对数字连接起来，并存储在输入的 Signature 字段。
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[inID].Signature = signature
	}
}

func (tx *Transcation) Verify(prevTXs map[string]Transcation) bool {
	fmt.Printf("开始验证 tx:%x\n",tx.TxHash)
	if tx.IsCoinbaseTransaction() {
		fmt.Printf("是CoinbaseTransaction tx:%x\n",tx.TxHash)
		return true
	}
	txCopy := tx.TrimmedCopy()
	//Next, we’ll need the same curve that is used to generate key pairs:
	//椭圆曲线
	curve := elliptic.P256()

	//检查每个在input中的签名
	for inID, vin := range tx.Vins {

		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.VoutIndex].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil
		//将TXI的Signature和PubKey中的数据进行“拆包”，用于crypto/ecdsa库进行验证使用
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])
		//ecdsa-椭圆曲线数字签名算法
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			fmt.Printf("开始验证 tx失败:%x\n",tx.TxHash)
			return false
		}
	}
	fmt.Printf("开始验证 tx成功:%x\n",tx.TxHash)
	return true
}

func (tx *Transcation) Hash() [] byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(tx.Serialize())
	return hash[:]
}

func (tx *Transcation) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

/**
复制一份交易，但input里的公钥和签名设置为null
 */
func (tx *Transcation) TrimmedCopy() Transcation {
	var inputs []*TXInput
	var outputs []*TXOutput
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.TxHash, vin.VoutIndex, nil, nil})
	}
	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transcation{tx.TxHash, inputs, outputs}
	return txCopy
}
