package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"time"
)

type Block struct {
	PrevHash     [32]byte
	Transactions []*Tx
	Timestamp    int64
	Nonce        int
}

func NewBlock(prevHash [32]byte, transactions []*Tx) *Block {
	b := &Block{
		Transactions: transactions,
		PrevHash:     prevHash,
		Timestamp:    0,
		Nonce:        0,
	}

	pow := NewProofOfWork(b)
	nonce := pow.Proof()

	b.Nonce = nonce
	b.Timestamp = time.Now().UnixNano()

	return b
}

func CreateGenesisBlock() *Block {
	return NewBlock([32]byte{}, []*Tx{})
}

func (b *Block) ToBytes() ([]byte, error) {
	return json.Marshal(struct {
		PrevHash     [32]byte `json:"prev_hash"`
		Transactions []*Tx    `json:"transactions"`
		Nonce        int      `json:"nonce"`
		Timestamp    int64    `json:"timestamp"`
	}{
		PrevHash:     b.PrevHash,
		Transactions: b.Transactions,
		Nonce:        b.Nonce,
		Timestamp:    b.Timestamp,
	})
}

func (b *Block) Hash() [32]byte {
	data, err := b.ToBytes()
	if err != nil {
		log.Panic(err)
	}
	return sha256.Sum256(data)
}
