package blockchain

import (
	"fmt"
	"testing"
	"time"
)

func TestPoW_ValidateReturnsTrue(t *testing.T) {
	pow := NewProofOfWork(&Block{
		Transactions: []*Tx{},
		PrevHash:     fmt.Sprintf("%x", [32]byte{}),
		Timestamp:    time.Now().UnixNano(),
		Nonce:        0,
	})

	nonce := pow.Run()

	if !pow.IsValid(nonce) {
		t.Error("IsValid should equal to true")
	}
}

func TestPoW_ValidateReturnsFalse(t *testing.T) {
	pow := NewProofOfWork(&Block{
		Transactions: []*Tx{},
		PrevHash:     fmt.Sprintf("%x", [32]byte{}),
		Timestamp:    time.Now().UnixNano(),
		Nonce:        0,
	})

	if pow.IsValid(0) {
		t.Error("IsValid should equal to false")
	}
}
