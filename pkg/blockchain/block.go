package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"time"
)

type Block struct {
	PrevHash     []byte
	Transactions []*Tx
	Timestamp    int64
	Nonce        int
}

func createBlock(prevHash []byte, tx *Tx) *Block {
	b := &Block{
		Transactions: []*Tx{tx},
		PrevHash:     prevHash,
		Timestamp:    time.Now().UnixNano(),
		Nonce:        0,
	}

	pow := NewProofOfWork(b)
	nonce := pow.Mine()

	b.Nonce = nonce

	return b
}

func CreateGenesisBlock() *Block {
	return createBlock([]byte{}, &Tx{})
}

func (b *Block) ToBytes() ([]byte, error) {
	return json.Marshal(struct {
		PrevHash     []byte `json:"prev_hash"`
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

func (b *Block) Hash() [32]byte {
	data, err := b.ToBytes()
	if err != nil {
		log.Panic(err)
	}
	return sha256.Sum256(data)
}
