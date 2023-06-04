package blockchain

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
)

type Tx struct {
	senderAddr    string
	recipientAddr string
	value         float32
}

func NewTransaction(sender string, recipient string, value float32) *Tx {
	return &Tx{
		senderAddr:    sender,
		recipientAddr: recipient,
		value:         value,
	}
}

func (tx *Tx) isCoinbase() bool {
	return tx.senderAddr == MINING_SENDER
}

func (tx *Tx) GenerateSignature(privateKey *ecdsa.PrivateKey) *utils.Signature {
	t, _ := tx.MarshalJSON()
	h := sha256.Sum256(t)
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    tx.senderAddr,
		Recipient: tx.recipientAddr,
		Value:     tx.value,
	})
}
