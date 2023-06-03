package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"log"
)

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

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32, s *utils.Signature, senderPublicKey *ecdsa.PublicKey) bool {
	tx := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		chain.TxPool = append(chain.TxPool, tx)
		return true
	}

	if chain.VerifyTxSignature(senderPublicKey, s, tx) {
		chain.TxPool = append(chain.TxPool, tx)
		return true
	} else {
		log.Panicln("ERROR: Failed to very tx")
	}
	return false
}

func (chain *Blockchain) VerifyTxSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, tx *Tx) bool {
	t, _ := tx.ToBytes()
	hash := sha256.Sum256(t)
	return ecdsa.Verify(senderPublicKey, hash[:], s.R, s.S)
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
	chain.AddTransaction(MINING_SENDER, chain.Address, MINING_REWARD, nil, nil)
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
