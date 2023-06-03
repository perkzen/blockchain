package blockchain

type Blockchain struct {
	Blocks  []*Block
	TxPool  []*Tx
	Address string
}

func (chain *Blockchain) AddBlock() {
	prevHash := chain.lastBlock().Hash()
	chain.Blocks = append(chain.Blocks, NewBlock(prevHash, chain.TxPool))
	chain.TxPool = []*Tx{}
}

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	chain.TxPool = append(chain.TxPool, NewTransaction(sender, recipient, value))
}

func (chain *Blockchain) CalculateTotalAmount(addr string) float32 {
	var total float32 = 0.0
	for _, tx := range chain.TxPool {
		if addr == tx.recipientAddr {
			total += tx.value
		}
		if addr == tx.senderAddr {
			total -= tx.value
		}
	}
	return total
}

func (chain *Blockchain) lastBlock() *Block {
	return chain.Blocks[len(chain.Blocks)-1]
}

func (chain *Blockchain) Mining() bool {
	chain.AddTransaction(MINING_SENDER, chain.Address, MINING_REWARD)
	chain.AddBlock()
	return true
}

func InitBlockchain(addr string) *Blockchain {
	return &Blockchain{
		Blocks:  []*Block{CreateGenesisBlock()},
		TxPool:  []*Tx{},
		Address: addr,
	}
}
