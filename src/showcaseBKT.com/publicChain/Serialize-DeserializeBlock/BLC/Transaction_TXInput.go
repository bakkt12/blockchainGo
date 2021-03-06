package BLC

import "bytes"

/**
交易输入由上个交易输出点、交易解锁脚本及序列号组成，
其中上个交易输出点包含两个元素，一个是上一个交易的哈希值，
另一个是上一个交易输出的索引号，由这两个元素便可确定唯一的UTXO
**/
/** An outpoint - a combination of a transaction hash and an index n into its vout */
type TXInput struct {
	TxHash    []byte //1.交易id
	VoutIndex int    //2.存储Txoutput在TXOutput中的索引
	//3 交易解锁脚本 btc中解锁脚本由签名和公钥组成
	//解锁脚本(scriptSig)   <sign> <PubK>
	//--->包含付款人对本次交易的签名(<sig>)和付款人公钥(<PubK(A)>)。
	Signature []byte //数字签名
	PublicKey []byte // 公钥  原生的，没有加密的
}

//当前消费的是谁的钱
func (in *TXInput) UnlockRipedm160Hash(ripemd160hash []byte) bool {
	publicKey := Ripemd160Hash(in.PublicKey)
	return bytes.Compare(publicKey, ripemd160hash) == 0
}
