package BLC

//转帐
func (cli *CLI) Send(from []string, to []string, amount []string) {

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amount)
}
