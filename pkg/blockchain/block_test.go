package blockchain

import (
	"fmt"
	"testing"
)

func TestBlock_Hash(t *testing.T) {
	b := &Block{
		Transactions: []*Tx{},
		PrevHash:     fmt.Sprintf("%x", [32]byte{}),
		Timestamp:    0,
		Nonce:        0,
	}

	a := &Block{
		Transactions: []*Tx{
			CoinbaseTx("a"),
			CoinbaseTx("b"),
		},
		PrevHash:  fmt.Sprintf("%x", [32]byte{}),
		Timestamp: 0,
		Nonce:     0,
	}

	c := &Block{
		Transactions: []*Tx{
			CoinbaseTx("a"),
			CoinbaseTx("b"),
		},
		PrevHash:  fmt.Sprintf("%x", [32]byte{}),
		Timestamp: 0,
		Nonce:     0,
	}

	d := &Block{
		Transactions: []*Tx{
			CoinbaseTx("b"),
			CoinbaseTx("a"),
		},
		PrevHash:  fmt.Sprintf("%x", [32]byte{}),
		Timestamp: 0,
		Nonce:     0,
	}

	hash := b.Hash()

	if hash == [32]byte{} {
		t.Error("Hash should not be empty")
	}

	if hash != b.Hash() {
		t.Error("Hash should be the same")
	}

	if a.Hash() != c.Hash() {
		t.Error("Hash should be the same")
	}

	if a.Hash() == d.Hash() {
		t.Error("Hash should be different")
	}
}
