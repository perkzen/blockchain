package blockchain

import (
	"fmt"
	"testing"
)

func TestBlock_Hash(t *testing.T) {

	b := &Block{
		Transactions: []*Tx{},
		PrevHash:     fmt.Sprintf("%x", []byte{}),
		Timestamp:    0,
		Nonce:        0,
	}

	hash := b.Hash()

	if hash == fmt.Sprintf("%x", []byte{}) {
		t.Error("Hash should not be empty")
	}
}
