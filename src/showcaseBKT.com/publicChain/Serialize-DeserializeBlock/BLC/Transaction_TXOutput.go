package BLC

import "bytes"

/**
交易输出由比特币数量、锁定脚本组成，
其中比特币数量表明了该输出包含的比特币数量，
锁定脚本对UTXO上了“锁”，谁能提供有效的解锁脚本，谁就能花费该UTXO。
 */
type TXOutput struct {
	Value      int64  //比特币数量
	Ripemd160Hash []byte //这里简单
	//ScriptPubKey []byte //锁定脚本scriptPubKey
	// 锁定脚本  <sign> <PubK> DUP HASH160  <PubkHash> EQUALVERIFY  CHECKSIG
	//1.<sign>签名 放到栈顶
	//2  <PubK> 公钥放到栈顶 原始的未加密
	//3  DUP 复制栈顶的值，也就是<PubK>
	//4  HASH160 ，对栈顶的值做 RIPEMD160(SHA256(PubK))
	//5 <PubkHash>  ,锁定脚 本的<PubkHash>放到栈顶
	//6 EQUALVERIFY，操作符是 把锁定脚 本的<PubkHash> 和计算RIPEMD160(SHA256(PubK)) 比较，如果一致，则删除，继续验证
	//7 CHECKSIG，操作符签名<sign> 和公钥的<sign> 是否一致，如果匹配则在顶部显示TRUE
}

func NewTXOutput(value int64, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

//从地址中解析出公钥 ，再设置回txoutput
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash [1 : len(pubKeyHash)-4]
	out.Ripemd160Hash = pubKeyHash
}

//address  : version(1) + 20位 +4位验证(4) -》Base58Decode
//从地址中解析出公钥
func (out *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	publicKeyHash := Base58Encode([]byte(address))
	RIPEMD160Hasher := publicKeyHash[1 : len(publicKeyHash)-4]
	return bytes.Compare(RIPEMD160Hasher, out.Ripemd160Hash) == 0
}
