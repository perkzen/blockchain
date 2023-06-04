package wallet

import (
	"testing"
)

func TestWallet_BlockchainAddress(t *testing.T) {
	w := NewWallet()
	if len(w.BlockchainAddress()) != 34 {
		t.Error("Blockchain address should be 34 characters long")
	}
}

func TestWallet_PrivateKey(t *testing.T) {
	w := NewWallet()
	if w.PrivateKey() == nil {
		t.Error("Private key should not be nil")
	}
}

func TestWallet_PublicKey(t *testing.T) {
	w := NewWallet()
	if w.PublicKey() == nil {
		t.Error("Public key should not be nil")
	}
}

func TestWallet_PrivateKeyStr(t *testing.T) {
	w := NewWallet()
	if len(w.PrivateKeyStr()) != 64 {
		t.Error("Private key string should be 64 characters long")
	}
}

func TestWallet_PublicKeyStr(t *testing.T) {
	w := NewWallet()
	if len(w.PublicKeyStr()) != 128 {
		t.Error("Public key string should be 128 characters long")
	}
}
