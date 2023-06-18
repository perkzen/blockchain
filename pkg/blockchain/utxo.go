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

// UTXO represents unspent transaction outputs.
// outputs is a list of all outputs.
// isSpent is a map of all outputs that have been spent.
type UTXO struct {
	outputs []Output
	isSpent map[string]bool
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
		outputs: []Output{},
		isSpent: make(map[string]bool),
	}
}

func (utxo *UTXO) AddOutput(output Output) {
	utxo.outputs = append(utxo.outputs, output)
	utxo.isSpent[output.Hash()] = false
}

func (utxo *UTXO) SpendOutput(hash string) {
	utxo.isSpent[hash] = true
}

func (utxo *UTXO) FindSpendableOutputs(address string, amount float32) (float32, []Output) {
	var unspentOut []Output
	var acc float32 = 0.0

Work:
	for _, out := range utxo.outputs {
		if out.TxOutput.CanBeUnlocked(address) {
			if !utxo.isSpent[out.Hash()] {
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

	for _, out := range utxo.outputs {
		if out.TxOutput.CanBeUnlocked(address) {
			if !utxo.isSpent[out.Hash()] {
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
