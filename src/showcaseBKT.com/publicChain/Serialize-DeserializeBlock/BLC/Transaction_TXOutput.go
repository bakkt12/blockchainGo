package BLC


/**
交易输出由比特币数量、锁定脚本组成，
其中比特币数量表明了该输出包含的比特币数量，
锁定脚本对UTXO上了“锁”，谁能提供有效的解锁脚本，谁就能花费该UTXO。
 */
type TXOutput struct {
	Value      int64    //比特币数量
	ScriptPubKey string //锁定脚本scriptPubKey
}

func (out *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	return out.ScriptPubKey == address
}