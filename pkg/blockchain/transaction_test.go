package blockchain

import (
	"blockchain/pkg/wallet"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	tx := NewTransaction("sender", "recipient", 0.1)
	if tx.senderAddr != "sender" {
		t.Error("Sender address should be set")
	}
	if tx.recipientAddr != "recipient" {
		t.Error("Recipient address should be set")
	}
	if tx.value != 0.1 {
		t.Error("Value should be set")
	}
}

func TestTx_isCoinbase(t *testing.T) {
	tx := NewTransaction("sender", "recipient", 0.1)
	if tx.isCoinbase() {
		t.Error("Coinbase should be false")
	}
	tx.senderAddr = MINING_SENDER
	if !tx.isCoinbase() {
		t.Error("Coinbase should be true")
	}
}

func TestTx_GenerateSignature(t *testing.T) {
	w := wallet.NewWallet()
	tx := NewTransaction("sender", "recipient", 0.1)
	signature := tx.GenerateSignature(w.PrivateKey())
	if signature == nil {
		t.Error("Signature should not be nil")
	}
}
