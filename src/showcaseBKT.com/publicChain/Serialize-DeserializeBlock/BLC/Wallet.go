package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	//1 私钥(用到椭圆曲线)
	PrivateKey ecdsa.PrivateKey

	//2 公钥
	PublicKey []byte
}

func IsValidForAddress(address []byte) bool {
	/**
	Version  Public key hash                           Checksum
	00       62E907B15CBF27D5425399EBF6F0FB50EBB88F18  C29B7D93

	 */
	fmt.Printf("开始检查地址: %s\n", address)
	version_public_checkSumBytes := Base58Decode(address)
	//fmt.Printf("version_public_checkSumBytes: %c ,len:%d \n", (version_public_checkSumBytes), len(version_public_checkSumBytes))
	//fmt.Println(address)
	//fmt.Println(version_public_checkSumBytes)
	//最后面4个字节（检查字节）,从（总长度-4 ）到未尾
	checkSumBytes := version_public_checkSumBytes[len(version_public_checkSumBytes)-addressChecksumLen:]
	//前面21个字节（version+ ripem160)，从0 到未尾倒数第4
	versionRipemd160Bytes := version_public_checkSumBytes[:len(version_public_checkSumBytes)-addressChecksumLen]

	//fmt.Printf("%s\n",checkSumBytes)
	//fmt.Printf("%s\n",versionRipemd160Bytes)

	checkBytes := checksum(versionRipemd160Bytes)

	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		fmt.Println("检查地址合法......")
		return true
	}
	fmt.Println("检查地址【不合法】  请再次确定检查......")
	return false;
}
func (w *Wallet) Getaddress() []byte {
	//1.使用 RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
	ripemd160Hash := Ripemd160Hash(w.PublicKey)

	//2.给哈希加上地址生成算法版本的前缀
	versionRipemd160Hash := append([]byte{version}, ripemd160Hash...)

	//3.对于第二步生成的结果，使用 SHA256(SHA256(payload)) 再哈希，计算校验和。校验和是结果哈希的前四个字节。
	checkSumBytes := checksum(versionRipemd160Hash)

	//4将校验和附加到 version+PubKeyHash 的组合中。

	// 5使用 Base58 对 version+PubKeyHash+checksum 组合进行编码
	bytes := append(versionRipemd160Hash, checkSumBytes ...)

	//就可以得到一个真实的比特币地址，你甚至可以在 blockchain.info 查看它的余额
	return Base58Encode(bytes)
}

func checksum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}

func Ripemd160Hash(publicKey []byte) []byte {
	//1.先将 publickey  SHA256
	sha256 := sha256.New()
	sha256.Write([]byte(publicKey))
	bytes := sha256.Sum(nil)

	//2 再做 RIP160
	ripemd160 := ripemd160.New()
	ripemd160.Write(bytes)
	return ripemd160.Sum(nil)
}

func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

//通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes() ...)

	return *private, pubKey
}
