package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// TxInput are references to previous Outputs.
// ID is the ID of the transaction that contains the output.
// OutIdx is the index of the output in the transaction.
type TxInput struct {
	ID        string `json:"id"`
	OutIdx    int    `json:"out_idx"`
	Signature string `json:"signature"`
}

// TxOutput represents a transaction output.
// Value is the amount of coins locked.
// PublicKey value needed to unlock the coins.
type TxOutput struct {
	Value     float32 `json:"value"`
	PublicKey string  `json:"public_key"`
}

type Tx struct {
	ID        string     `json:"id"`
	TxInputs  []TxInput  `json:"tx_inputs"`
	TxOutputs []TxOutput `json:"tx_outputs"`
	Timestamp int64      `json:"timestamp"`
}

func NewTransaction(sender string, recipient string, value float32, chain *Blockchain) *Tx {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.UTXO.FindSpendableOutputs(sender, value)

	if acc < value {
		log.Panic("ERROR: Insufficient funds.")
	}

	for _, out := range validOutputs {
		input := TxInput{ID: out.TxID, OutIdx: out.Idx, Signature: sender}
		inputs = append(inputs, input)
		chain.UTXO.SpendOutput(out.Hash())
	}

	txOut := TxOutput{Value: value, PublicKey: recipient}
	outputs = append(outputs, txOut)

	if acc > value {
		outputs = append(outputs, TxOutput{Value: acc - value, PublicKey: sender})
	}

	tx := Tx{TxInputs: inputs, TxOutputs: outputs, Timestamp: time.Now().UnixNano()}
	tx.setID()

	chain.UTXO.AddOutput(NewOutput(txOut, 0, tx.ID, tx.Timestamp))
	if acc > value {
		chain.UTXO.AddOutput(NewOutput(outputs[1], 1, tx.ID, tx.Timestamp))
	}

	return &tx
}

func (tx *Tx) isCoinbase() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].ID) == 0 && tx.TxInputs[0].OutIdx == -1
}

func (in *TxInput) CanUnlock(addr string) bool {
	return in.Signature == addr
}

func (out *TxOutput) CanBeUnlocked(addr string) bool {
	return out.PublicKey == addr
}

func CoinbaseTx(chain *Blockchain) *Tx {
	txIn := TxInput{ID: "", OutIdx: -1, Signature: chain.Address}
	txOut := TxOutput{Value: COINBASE_REWARD, PublicKey: chain.Address}
	cbTx := &Tx{
		ID:       fmt.Sprintf("%x", [32]byte{}),
		TxInputs: []TxInput{txIn}, TxOutputs: []TxOutput{txOut},
		Timestamp: time.Now().UnixNano(),
	}
	cbTx.setID()
	chain.UTXO.AddOutput(NewOutput(txOut, 0, cbTx.ID, cbTx.Timestamp))
	return cbTx
}

func (tx *Tx) GenerateSignature(privateKey *ecdsa.PrivateKey) *utils.Signature {
	t, _ := tx.MarshalJSON()
	h := sha256.Sum256(t)
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (tx *Tx) setID() {
	t, _ := tx.MarshalJSON()
	tx.ID = fmt.Sprintf("%x", sha256.Sum256(t[:]))
}

func (out *TxOutput) Hash() string {
	t, _ := out.MarshalJSON()
	return fmt.Sprintf("%x", sha256.Sum256(t[:]))
}

func (out *TxOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value     float32 `json:"value"`
		PublicKey string  `json:"public_key"`
	}{
		Value:     out.Value,
		PublicKey: out.PublicKey,
	})
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string     `json:"id"`
		TxInputs  []TxInput  `json:"tx_inputs"`
		TxOutputs []TxOutput `json:"tx_outputs"`
		Timestamp int64      `json:"timestamp"`
	}{
		ID:        tx.ID,
		TxInputs:  tx.TxInputs,
		TxOutputs: tx.TxOutputs,
		Timestamp: tx.Timestamp,
	})
}
