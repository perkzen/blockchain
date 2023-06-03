package blockchain

type Blockchain struct {
	Blocks []*Block
}

func (chain *Blockchain) AddBlock(tx *Tx) {
	prevHash := chain.lastBlock().PrevHash
	chain.Blocks = append(chain.Blocks, createBlock(prevHash, tx))
}

func (chain *Blockchain) lastBlock() *Block {
	return chain.Blocks[len(chain.Blocks)-1]
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{CreateGenesisBlock()}}
}
