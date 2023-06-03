package blockchain

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func isValidHash(hash [32]byte) bool {
	guess := fmt.Sprintf("%x", hash)
	zeros := strings.Repeat("0", DIFFICULTY)
	return guess[:DIFFICULTY] == zeros
}

func TestPoW_Mine(t *testing.T) {
	pow := NewProofOfWork(&Block{
		Data:      []byte("test"),
		PrevHash:  []byte{},
		Timestamp: time.Now().UnixNano(),
		Nonce:     0,
	})

	hash, _ := pow.Mine()

	if !isValidHash(hash) {
		t.Error("Invalid hash format")
	}
}

func TestPoW_ValidateReturnsTrue(t *testing.T) {
	pow := NewProofOfWork(&Block{
		Data:      []byte("test"),
		PrevHash:  []byte{},
		Timestamp: time.Now().UnixNano(),
		Nonce:     0,
	})

	_, nonce := pow.Mine()

	if !pow.Validate(nonce) {
		t.Error("Validate should equal to true")
	}
}

func TestPoW_ValidateReturnsFalse(t *testing.T) {
	pow := NewProofOfWork(&Block{
		Data:      []byte("test"),
		PrevHash:  []byte{},
		Timestamp: time.Now().UnixNano(),
		Nonce:     0,
	})

	if pow.Validate(0) {
		t.Error("Validate should equal to false")
	}
}
