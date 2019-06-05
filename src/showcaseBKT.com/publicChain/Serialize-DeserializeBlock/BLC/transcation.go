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
	ID   []byte     // 交易id
	Vin  []TXInput  // 交易输入
	Vout []TXOutput //交易输出
	//	Lock_time int64 未用到
}

// 1.判断当前交易是否是 coinbase tx
func (tx *Transcation) IsCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].VoutIndex == -1 && len(tx.Vin[0].Txid) == 0
}

func (transcation *Transcation) printfTranscation() {

	fmt.Printf("\t ####txid:		%x\n", transcation.ID)
	for _, in := range transcation.Vin {
		fmt.Println("\t-------vinput----------")
		fmt.Printf("\tvin txid        :%x\n", in.Txid)
		fmt.Printf("\tvin voutIndex   :%d\n", in.VoutIndex)
		fmt.Printf("\tvin ScriptPubKey:%s\n", in.ScriptPubKey)
	}
//	fmt.Println("")
	for _, out := range transcation.Vout {
		fmt.Println("\t--------vout---------")
		fmt.Printf("\tvout amount      :%d\n", out.CAmount)
		fmt.Printf("\tvout ScriptPubKey:%s\n", out.ScriptPubKey)
	}

}

//创建一个新的coinbase交易
func NewCoinbaseTx(to, data string) *Transcation {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	//创建创世的输入
	txin := TXInput{[]byte{}, -1, data}
	//创建输出
	txout := TXOutput{subsildy, to}
	//创建交易
	tx := Transcation{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

//建立交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain,noPackageTxs []*Transcation) *Transcation {

	fmt.Printf("Create a new UTXOTransaction..from:%s->to:%s \n", from, to)
	fmt.Println("start printf noPackageTxs")
	for _,tx := range  noPackageTxs{
		tx.printfTranscation()
	}
	fmt.Println("end printf noPackageTxs")
	//输入
	var inputs []TXInput
	//输出
	var outputs []TXOutput

	//1.找到有效的可用的交易输出数据模型
	//查询出未花费的输出  (int,map[string][]int)
	acc, validOutputs := bc.FindSpendableOutputs(from, amount,noPackageTxs)

	if acc < amount {
		fmt.Printf("acc:%d  amount:%d \n",acc,amount)
		log.Panic("ERROR; Not enough funds!")
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//消费掉未消费的out，用vout来创建 vin
		for _, voutIndex := range outs {
			//创建一个输入
			input := TXInput{txID, voutIndex, from}
			inputs = append(inputs, input)
		}
	}

	//建立输出转帐
	output := TXOutput{amount, to}
	outputs = append(outputs, output)
	//建立输出，找零
	output = TXOutput{acc - amount, from}
	outputs = append(outputs, output)

	//创建交易
	tx := Transcation{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

//设置交易 hash，将 tx交易序列化之后字节数组生成256Hash
func (tx *Transcation) SetID() {
	var encoded bytes.Buffer;
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

/**
交易输入由上个交易输出点、交易解锁脚本及序列号组成，
其中上个交易输出点包含两个元素，一个是上一个交易的哈希值，
另一个是上一个交易输出的索引号，由这两个元素便可确定唯一的UTXO
**/
/** An outpoint - a combination of a transaction hash and an index n into its vout */
type TXInput struct {
	Txid         []byte //1.交易id
	VoutIndex    int    //2.存储Txoutput在TXOutput中的索引
	ScriptPubKey string //3 交易解锁脚本 btc中解锁脚本由签名和公钥组成
}

/**
交易输出由比特币数量、锁定脚本组成，
其中比特币数量表明了该输出包含的比特币数量，
锁定脚本对UTXO上了“锁”，谁能提供有效的解锁脚本，谁就能花费该UTXO。
 */
type TXOutput struct {
	CAmount      int    //比特币数量
	ScriptPubKey string //锁定脚本scriptPubKey
}

//简单的帐号地址，是否能够解锁帐号
func (in *TXInput) CanUnlockOutputWith(unlockData string) bool {
	return in.ScriptPubKey == unlockData
}

func (out *TXOutput) CanBeUnlockedWith(unlockData string) bool {
	return out.ScriptPubKey == unlockData

}
