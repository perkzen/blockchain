package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// Output represents a transaction output.
// Idx is the index of the output in the transaction.
// TxID is the ID of the transaction that contains the output.
// Timestamp is the transaction timestamp.
type Output struct {
	TxOutput  TxOutput
	Idx       int
	TxID      string
	Timestamp int64
}

// UTXO represents unspent transaction Outputs.
// Outputs is a list of all Outputs.
// IsSpent is a map of all Outputs that have been spent.
type UTXO struct {
	Outputs []Output
	IsSpent map[string]bool
}

func NewOutput(txOutput TxOutput, idx int, txID string, timestamp int64) Output {
	return Output{
		TxOutput:  txOutput,
		Idx:       idx,
		TxID:      txID,
		Timestamp: timestamp,
	}
}

func NewUTXO() *UTXO {
	return &UTXO{
		Outputs: []Output{},
		IsSpent: make(map[string]bool),
	}
}

func (utxo *UTXO) AddOutput(output Output) {
	utxo.Outputs = append(utxo.Outputs, output)
	utxo.IsSpent[output.Hash()] = false
}

func (utxo *UTXO) SpendOutput(hash string) {
	utxo.IsSpent[hash] = true
}

func (utxo *UTXO) FindSpendableOutputs(address string, amount float32) (float32, []Output) {
	var unspentOut []Output
	var acc float32 = 0.0

Work:
	for _, out := range utxo.Outputs {
		if out.TxOutput.CanBeUnlocked(address) {
			if !utxo.IsSpent[out.Hash()] {
				acc += out.TxOutput.Value
				unspentOut = append(unspentOut, out)
				if acc >= amount {
					break Work
				}
			}
		}
	}

	return acc, unspentOut
}

func (utxo *UTXO) GetBalance(address string) float32 {
	var balance float32 = 0.0

	for _, out := range utxo.Outputs {
		if out.TxOutput.CanBeUnlocked(address) {
			if !utxo.IsSpent[out.Hash()] {
				balance += out.TxOutput.Value
			}
		}
	}

	return balance
}

func (out *Output) MarshallJSON() ([]byte, error) {
	return json.Marshal(struct {
		TxOutput struct {
			Value     float32 `json:"value"`
			PublicKey string  `json:"public_key"`
		} `json:"tx_output"`
		Idx       int    `json:"idx"`
		TxID      string `json:"txID"`
		Timestamp int64  `json:"timestamp"`
	}{
		TxOutput:  out.TxOutput,
		Idx:       out.Idx,
		TxID:      out.TxID,
		Timestamp: out.Timestamp,
	})
}

func (out *Output) Hash() string {
	o, _ := json.Marshal(out)
	return fmt.Sprintf("%x", sha256.Sum256(o[:]))
}
