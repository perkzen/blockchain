package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"log"
)

type Blockchain struct {
	blocks   []*Block
	txPool   []*Tx
	address  string
	port     uint16
	name     string
	currency string
}

func (chain *Blockchain) AddBlock() {
	prevHash := chain.lastBlock().Hash()
	chain.blocks = append(chain.blocks, NewBlock(prevHash, chain.txPool))
	chain.txPool = []*Tx{}
}

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32, s *utils.Signature, senderPublicKey *ecdsa.PublicKey) bool {
	tx := NewTransaction(sender, recipient, value)

	if tx.isCoinbase() {
		chain.txPool = append(chain.txPool, tx)
		return true
	}

	if chain.VerifyTxSignature(senderPublicKey, s, tx) {
		if chain.CalculateTotalAmount(sender) < value {
			log.Printf("ERROR: Insufficient funds")
			return false
		}

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
	for _, block := range chain.blocks {
		for _, tx := range block.Transactions {
			if addr == tx.recipientAddr {
				total += tx.value
			}
			if addr == tx.senderAddr {
				total -= tx.value
			}
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
		blocks:   []*Block{CreateGenesisBlock()},
		txPool:   []*Tx{},
		address:  addr,
		port:     port,
		name:     BLOCKCHAIN_NAME,
		currency: BLOCKCHAIN_CURRENCY,
	}
}

func (chain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks   []*Block `json:"blocks"`
		Name     string   `json:"name"`
		Currency string   `json:"currency"`
	}{
		Blocks:   chain.blocks,
		Name:     chain.name,
		Currency: chain.currency,
	})
}
