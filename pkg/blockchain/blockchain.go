package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"log"
)

// Blockchain
// Address: Address of the node (miner)
type Blockchain struct {
	Blocks   []*Block `json:"blocks"`
	TxPool   []*Tx    `json:"tx_pool"`
	UTXO     *UTXO    `json:"utxo"`
	Address  string   `json:"address"`
	Port     uint16   `json:"port"`
	Name     string   `json:"name"`
	Currency string   `json:"currency"`
}

func (chain *Blockchain) AddBlock() *Block {
	prevHash := chain.lastBlock().Hash()
	txPool := append([]*Tx{CoinbaseTx(chain)}, chain.TxPool...)
	newBlock := NewBlock(prevHash, txPool)
	chain.Blocks = append(chain.Blocks, newBlock)
	chain.TxPool = []*Tx{}
	return newBlock
}

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPrivateKey *ecdsa.PrivateKey, senderPublicKey *ecdsa.PublicKey) bool {
	tx := NewTransaction(sender, recipient, value, chain)
	s := tx.GenerateSignature(senderPrivateKey)
	if tx.isCoinbase() {
		chain.TxPool = append(chain.TxPool, tx)
		return true
	}
	if chain.VerifyTxSignature(senderPublicKey, s, tx) {
		chain.TxPool = append(chain.TxPool, tx)
		return true
	} else {
		log.Panicln("ERROR: Failed to verify tx")
	}
	return false
}

func (chain *Blockchain) VerifyTxSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, tx *Tx) bool {
	t, _ := tx.MarshalJSON()
	hash := sha256.Sum256(t)
	return ecdsa.Verify(senderPublicKey, hash[:], s.R, s.S)
}

func (chain *Blockchain) lastBlock() *Block {
	return chain.Blocks[len(chain.Blocks)-1]
}

func (chain *Blockchain) CreateGenesisBlock() *Block {
	return NewBlock("GENESIS", []*Tx{CoinbaseTx(chain)})
}

func (chain *Blockchain) MineBlock() *Block {
	return chain.AddBlock()
}

func InitBlockchain(addr string, port uint16) *Blockchain {
	blockchain := &Blockchain{
		Blocks:   []*Block{},
		TxPool:   []*Tx{},
		UTXO:     NewUTXO(),
		Address:  addr,
		Port:     port,
		Name:     BLOCKCHAIN_NAME,
		Currency: BLOCKCHAIN_CURRENCY,
	}

	blockchain.Blocks = append(blockchain.Blocks, blockchain.CreateGenesisBlock())

	return blockchain
}

func (chain *Blockchain) Length() int {
	return len(chain.Blocks)
}

func (chain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks   []*Block `json:"blocks"`
		TxPool   []*Tx    `json:"tx_pool"`
		UTXO     *UTXO    `json:"utxo"`
		Name     string   `json:"name"`
		Currency string   `json:"currency"`
	}{
		Blocks:   chain.Blocks,
		TxPool:   chain.TxPool,
		UTXO:     chain.UTXO,
		Name:     chain.Name,
		Currency: chain.Currency,
	})
}
