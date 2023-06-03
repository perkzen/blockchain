package blockchain

type Blockchain struct {
	Blocks []*Block
	TxPool []*Tx
}

func (chain *Blockchain) AddBlock() {
	prevHash := chain.lastBlock().Hash()
	chain.Blocks = append(chain.Blocks, NewBlock(prevHash, chain.TxPool))
	chain.TxPool = []*Tx{}
}

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	chain.TxPool = append(chain.TxPool, NewTransaction(sender, recipient, value))
}

func (chain *Blockchain) lastBlock() *Block {
	return chain.Blocks[len(chain.Blocks)-1]
}

func InitBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*Block{CreateGenesisBlock()}, TxPool: []*Tx{}}
}
