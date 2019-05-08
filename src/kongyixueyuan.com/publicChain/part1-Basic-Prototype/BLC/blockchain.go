package BLC

type Blockchain struct {
	Blocks []*Block //存储有序的区块
}

//新增区块
func (blockchain *Blockchain) AddBlock(data string) {
	//1.创新区块
	preBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	newBlock := NewBlock(data, preBlock.Hash)

	//2 加到blockchain里面
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

//创建一个带有创或区块的区块键
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
