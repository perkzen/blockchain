package wallet

import (
	"blockchain/pkg/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
)

type Tx struct {
	senderPublicKey  *ecdsa.PublicKey
	senderPrivateKey *ecdsa.PrivateKey
	senderAddr       string
	recipientAddr    string
	value            float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Tx {
	return &Tx{
		senderPrivateKey: privateKey,
		senderPublicKey:  publicKey,
		senderAddr:       sender,
		recipientAddr:    recipient,
		value:            value,
	}
}

func (tx *Tx) GenerateSignature() *utils.Signature {
	t, _ := tx.MarshalJSON()
	h := sha256.Sum256(t)
	r, s, _ := ecdsa.Sign(rand.Reader, tx.senderPrivateKey, h[:])
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
