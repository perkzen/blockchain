package blockchain

import "encoding/json"

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

func (tx *Tx) ToBytes() ([]byte, error) {
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
