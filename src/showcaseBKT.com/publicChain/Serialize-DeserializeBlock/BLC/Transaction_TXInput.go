package BLC

/**
交易输入由上个交易输出点、交易解锁脚本及序列号组成，
其中上个交易输出点包含两个元素，一个是上一个交易的哈希值，
另一个是上一个交易输出的索引号，由这两个元素便可确定唯一的UTXO
**/
/** An outpoint - a combination of a transaction hash and an index n into its vout */
type TXInput struct {
	TxHash          []byte //1.交易id
	VoutIndex    int    //2.存储Txoutput在TXOutput中的索引
	ScriptSig string //3 交易解锁脚本 btc中解锁脚本由签名和公钥组成
}

//简单的帐号地址，是否能够解锁帐号
func (in *TXInput) UnLockWithAddress(address string) bool {
	return in.ScriptSig == address
}