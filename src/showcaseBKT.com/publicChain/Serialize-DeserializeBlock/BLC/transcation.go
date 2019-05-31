package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsildy = 50

type Transcation struct {
	//Vesion string 未用到
	ID   []byte     // 交易id
	Vin  []TXInput  // 交易输入
	Vout []TXOutput //交易输出
	//	Lock_time int64 未用到
}

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
