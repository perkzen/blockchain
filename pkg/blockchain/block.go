package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Block struct {
	PrevHash     string `json:"prev_hash"`
	Transactions []*Tx  `json:"transactions"`
	Timestamp    int64  `json:"timestamp"`
	Nonce        int    `json:"nonce"`
}

func NewBlock(prevHash string, transactions []*Tx) *Block {
	b := &Block{
		Transactions: transactions,
		PrevHash:     prevHash,
		Timestamp:    time.Now().UnixNano(),
		Nonce:        0,
	}

	return b
}

func (b *Block) setNonce(nonce int) {
	b.Nonce = nonce
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrevHash     string `json:"prev_hash"`
		Transactions []*Tx  `json:"transactions"`
		Nonce        int    `json:"nonce"`
		Timestamp    int64  `json:"timestamp"`
	}{
		PrevHash:     b.PrevHash,
		Transactions: b.Transactions,
		Nonce:        b.Nonce,
		Timestamp:    b.Timestamp,
	})
}

func (b *Block) Hash() string {
	data, err := b.MarshalJSON()
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
