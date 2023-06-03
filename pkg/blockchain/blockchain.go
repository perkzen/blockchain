package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"log"
)

type Blockchain struct {
	blocks  []*Block
	txPool  []*Tx
	address string
	port    uint16
}

func (chain *Blockchain) AddBlock() {
	prevHash := chain.lastBlock().Hash()
	chain.blocks = append(chain.blocks, NewBlock(prevHash, chain.txPool))
	chain.txPool = []*Tx{}
}

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32, s *utils.Signature, senderPublicKey *ecdsa.PublicKey) bool {
	tx := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		chain.txPool = append(chain.txPool, tx)
		return true
	}

	if chain.VerifyTxSignature(senderPublicKey, s, tx) {
		//if chain.CalculateTotalAmount(sender) < value {
		//	log.Printf("ERROR: Insufficient funds")
		//	return false
		//}

		chain.txPool = append(chain.txPool, tx)
		return true
	} else {
		log.Panicln("ERROR: Failed to very tx")
	}
	return false
}

func (chain *Blockchain) VerifyTxSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, tx *Tx) bool {
	t, _ := tx.MarshalJSON()
	hash := sha256.Sum256(t)
	return ecdsa.Verify(senderPublicKey, hash[:], s.R, s.S)
}

func (chain *Blockchain) CalculateTotalAmount(addr string) float32 {
	var total float32 = 0.0
	for _, tx := range chain.txPool {
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
	return chain.blocks[len(chain.blocks)-1]
}

func (chain *Blockchain) Mining() bool {
	chain.AddTransaction(MINING_SENDER, chain.address, MINING_REWARD, nil, nil)
	chain.AddBlock()
	return true
}

func InitBlockchain(addr string, port uint16) *Blockchain {
	return &Blockchain{
		blocks:  []*Block{CreateGenesisBlock()},
		txPool:  []*Tx{},
		address: addr,
		port:    port,
	}
}

func (chain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blocks"`
	}{
		Blocks: chain.blocks,
	})
}
