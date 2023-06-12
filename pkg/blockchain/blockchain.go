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
// Address: Address of the node (miner)
type Blockchain struct {
	Blocks   []*Block `json:"blocks"`
	TxPool   []*Tx    `json:"tx_pool"`
	Address  string   `json:"address"`
	Port     uint16   `json:"port"`
	Name     string   `json:"name"`
	Currency string   `json:"currency"`
}

func (chain *Blockchain) AddBlock() *Block {
	prevHash := chain.lastBlock().Hash()
	txPool := append([]*Tx{CoinbaseTx(chain.Address)}, chain.TxPool...)
	newBlock := NewBlock(prevHash, txPool)
	chain.Blocks = append(chain.Blocks, newBlock)
	chain.TxPool = []*Tx{}
	return newBlock
}

func (chain *Blockchain) AddTransaction(sender string, recipient string, value float32, s *utils.Signature, senderPublicKey *ecdsa.PublicKey) bool {
	tx := NewTransaction(sender, recipient, value, chain)
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

func CreateGenesisBlock(addr string) *Block {
	return NewBlock("GENESIS", []*Tx{CoinbaseTx(addr)})
}

func (chain *Blockchain) MineBlock() *Block {
	return chain.AddBlock()
}

func (chain *Blockchain) FindUnspentTxs(address string) []Tx {
	var utxo []Tx
	spentTXOs := make(map[string][]int)

	for _, block := range chain.Blocks {
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
		Blocks:   []*Block{CreateGenesisBlock(addr)},
		TxPool:   []*Tx{},
		Address:  addr,
		Port:     port,
		Name:     BLOCKCHAIN_NAME,
		Currency: BLOCKCHAIN_CURRENCY,
	}
}

func (chain *Blockchain) Length() int {
	return len(chain.Blocks)
}

func (chain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks   []*Block `json:"blocks"`
		TxPool   []*Tx    `json:"tx_pool"`
		Name     string   `json:"name"`
		Currency string   `json:"currency"`
	}{
		Blocks:   chain.Blocks,
		TxPool:   chain.TxPool,
		Name:     chain.Name,
		Currency: chain.Currency,
	})
}
