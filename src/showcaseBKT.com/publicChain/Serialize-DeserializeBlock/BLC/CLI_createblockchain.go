package BLC



//创建创世区块
func (cli *CLI) createGenesisBlockchain(genesis string) {

	CreateBlockchainWithGenesisBlock(genesis)
}