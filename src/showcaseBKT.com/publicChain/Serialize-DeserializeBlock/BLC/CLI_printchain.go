package BLC

func (cli *CLI) PrintChain() {

	blockchian := BlockchainObject()
	defer blockchian.DB.Close()

	blockchian.Printchain()
}
