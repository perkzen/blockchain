package blockchain

import (
	"testing"
	"time"
)

func TestPoW_ValidateReturnsTrue(t *testing.T) {
	pow := NewProofOfWork(&Block{
		Transactions: []*Tx{},
		PrevHash:     [32]byte{},
		Timestamp:    time.Now().UnixNano(),
		Nonce:        0,
	})

	nonce := pow.Proof()

	if !pow.Validate(nonce) {
		t.Error("Validate should equal to true")
	}
}

func TestPoW_ValidateReturnsFalse(t *testing.T) {
	pow := NewProofOfWork(&Block{
		Transactions: []*Tx{},
		PrevHash:     [32]byte{},
		Timestamp:    time.Now().UnixNano(),
		Nonce:        0,
	})

	if pow.Validate(0) {
		t.Error("Validate should equal to false")
	}
}
