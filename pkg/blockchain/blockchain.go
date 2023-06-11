package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
)

// Blockchain
// address: address of the node (miner)
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
	txPool := append([]*Tx{CoinbaseTx(chain.address)}, chain.txPool...)
	chain.blocks = append(chain.blocks, NewBlock(prevHash, txPool))
	chain.txPool = []*Tx{}
}

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32, s *utils.Signature, senderPublicKey *ecdsa.PublicKey) bool {
	tx := NewTransaction(sender, recipient, value, chain)
	if tx.isCoinbase() {
		chain.txPool = append(chain.txPool, tx)
		return true
	}
	if chain.VerifyTxSignature(senderPublicKey, s, tx) {
		chain.txPool = append(chain.txPool, tx)
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
	return chain.blocks[len(chain.blocks)-1]
}

func CreateGenesisBlock(addr string) *Block {
	return NewBlock([32]byte{}, []*Tx{CoinbaseTx(addr)})
}

func (chain *Blockchain) MineBlock() {
	chain.AddBlock()
}

func (chain *Blockchain) FindUnspentTxs(address string) []Tx {
	var utxo []Tx
	spentTXOs := make(map[string][]int)

	for _, block := range chain.blocks {
		for _, tx := range block.Transactions {
			txID := fmt.Sprintf("%x", tx.ID)
		Outputs:
			for outIdx, txOut := range tx.TxOutputs {
				if spentTXOs[txID] != nil {
					for _, spent := range spentTXOs[txID] {
						if spent == outIdx {
							continue Outputs
						}
					}
				}
				if txOut.CanBeUnlocked(address) {
					utxo = append(utxo, *tx)
				}
			}
			if !tx.isCoinbase() {
				for _, txIn := range tx.TxInputs {
					if txIn.CanUnlock(address) {
						inTxID := fmt.Sprintf("%x", txIn.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], txIn.OutIdx)
					}
				}
			}
		}
		if len(block.PrevHash[:]) == 0 {
			break
		}
	}
	return utxo
}

func (chain *Blockchain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	txs := chain.FindUnspentTxs(address)

	for _, tx := range txs {
		for _, out := range tx.TxOutputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (chain *Blockchain) FindSpendableOutputs(address string, amount float32) (float32, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTxs := chain.FindUnspentTxs(address)
	accumulated := float32(0.0)

Work:
	for _, tx := range unspentTxs {
		txID := fmt.Sprintf("%x", tx.ID)
		for outIdx, out := range tx.TxOutputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func InitBlockchain(addr string, port uint16) *Blockchain {
	return &Blockchain{
		blocks:   []*Block{CreateGenesisBlock(addr)},
		txPool:   []*Tx{},
		address:  addr,
		port:     port,
		name:     BLOCKCHAIN_NAME,
		currency: BLOCKCHAIN_CURRENCY,
	}
}

func (chain *Blockchain) Address() string {
	return chain.address
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
