package blockchain

import (
	"blockchain/pkg/wallet"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	chain := InitBlockchain(walletA.BlockchainAddress(), 3000)
	tx := NewTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 0.1, chain)

	if tx == nil {
		t.Error("Transaction should not be nil")
	}
	if tx.TxOutputs[0].Value != 0.1 {
		t.Error("Transaction value should be 0.1")
	}
	if tx.TxOutputs[0].PublicKey != walletB.BlockchainAddress() {
		t.Error("Transaction recipient should be walletB")
	}
	if tx.TxInputs[0].Signature != walletA.BlockchainAddress() {
		t.Error("Transaction sender should be walletA")
	}
}

func TestTx_isCoinbase(t *testing.T) {
	miner := wallet.NewWallet()
	chain := InitBlockchain(miner.BlockchainAddress(), 3000)
	ctx := CoinbaseTx(chain)

	if !ctx.isCoinbase() {
		t.Error("Transaction should be coinbase")
	}
}

func TestTx_GenerateSignature(t *testing.T) {
	w := wallet.NewWallet()
	chain := InitBlockchain(w.BlockchainAddress(), 3000)
	tx := NewTransaction(w.BlockchainAddress(), "recipient", 0.1, chain)
	signature := tx.GenerateSignature(w.PrivateKey())
	if signature == nil {
		t.Error("Signature should not be nil")
	}

	isValid := chain.VerifyTxSignature(w.PublicKey(), signature, tx)
	if !isValid {
		t.Error("Transaction signature should be valid")
	}

}
