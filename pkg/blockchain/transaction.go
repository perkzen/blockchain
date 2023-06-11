package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
)

// TxInput are references to previous outputs.
type TxInput struct {
	ID        []byte
	OutIdx    int
	Signature string
}

// TxOutput represents a transaction output.
// Value is the amount of coins locked.
// PublicKey value needed to unlock the coins.
type TxOutput struct {
	Value     float32
	PublicKey string
}

type Tx struct {
	ID        [32]byte
	TxInputs  []TxInput
	TxOutputs []TxOutput
}

func NewTransaction(sender string, recipient string, value float32, chain *Blockchain) *Tx {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(sender, value)

	if acc < value {
		log.Panic("ERROR: Insufficient funds.")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := TxInput{ID: txID, OutIdx: out, Signature: sender}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{Value: value, PublicKey: recipient})

	if acc > value {
		outputs = append(outputs, TxOutput{Value: acc - value, PublicKey: sender})
	}

	tx := Tx{TxInputs: inputs, TxOutputs: outputs}
	tx.setID()
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

func CoinbaseTx(receiverAddr string) *Tx {
	txIn := TxInput{ID: []byte{}, OutIdx: -1, Signature: receiverAddr}
	txOut := TxOutput{Value: COINBASE_REWARD, PublicKey: receiverAddr}
	cbTx := &Tx{ID: [32]byte{}, TxInputs: []TxInput{txIn}, TxOutputs: []TxOutput{txOut}}
	cbTx.setID()
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
	tx.ID = sha256.Sum256(t[:])
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string     `json:"id"`
		TxInputs  []TxInput  `json:"tx_inputs"`
		TxOutputs []TxOutput `json:"tx_outputs"`
	}{
		ID:        fmt.Sprintf("%x", tx.ID),
		TxInputs:  tx.TxInputs,
		TxOutputs: tx.TxOutputs,
	})
}
