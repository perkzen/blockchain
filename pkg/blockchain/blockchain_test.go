package blockchain

import (
	"blockchain/pkg/wallet"
	"fmt"
	"testing"
)

func TestCreateGenesisBlock(t *testing.T) {
	genesis := CreateGenesisBlock()
	genesisHash := fmt.Sprintf("%x", genesis.PrevHash)
	hash := fmt.Sprintf("%x", [32]byte{})
	if genesisHash != hash {
		t.Error("Hashes do not equal")
	}
}

func TestInitBlockchain(t *testing.T) {
	chain := InitBlockchain("", 3000)
	if len(chain.blocks) < 1 || len(chain.blocks) > 1 {
		t.Error("Chain should have only 1 block on initiation")
	}
}

func TestBlockchain_AddBlock(t *testing.T) {
	chain := InitBlockchain("", 3000)
	chain.AddBlock()
	if len(chain.blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
}

func TestBlockchain_AddTransaction(t *testing.T) {
	chain := InitBlockchain("", 3000)
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	tx := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	isAdd := chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, tx.GenerateSignature(), walletA.PublicKey())

	if !isAdd {
		t.Error("Tx should be added to tx pool")
	}

	if len(chain.txPool) < 1 {
		t.Error("Chain should have 1 transaction in pool")
	}
}

func TestBlockchain_ClearPool(t *testing.T) {
	chain := InitBlockchain("", 3000)
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	tx := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	isAdd := chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, tx.GenerateSignature(), walletA.PublicKey())

	if !isAdd {
		t.Error("Tx should be added to tx pool")
	}

	chain.AddBlock()
	if len(chain.blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
	if len(chain.txPool) >= 1 {
		t.Error("Transaction pool should be empty")
	}
}

func TestBlockchain_Mining(t *testing.T) {
	chain := InitBlockchain("", 3000)
	chain.Mining()
	lastBlock := chain.lastBlock()
	if len(chain.blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
	if len(chain.txPool) >= 1 {
		t.Error("Transaction pool should be empty")
	}
	if len(lastBlock.Transactions) != 1 {
		t.Error("Block should have 1 transaction in it")
	}
	if lastBlock.Transactions[0].value != MINING_REWARD {
		t.Error("Transaction value should equal mining reward")
	}
	if lastBlock.Transactions[0].senderAddr != MINING_SENDER {
		t.Error("Mining sender should equal THE BLOCKCHAIN")
	}
}

func TestBlockchain_CalculateTotalAmount(t *testing.T) {
	chain := InitBlockchain("", 3000)
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	chain.Mining()

	tx1 := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	chain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, tx1.GenerateSignature(), walletA.PublicKey())

	tx2 := wallet.NewTransaction(walletB.PrivateKey(), walletB.PublicKey(), walletB.BlockchainAddress(), walletA.BlockchainAddress(), 1.0)
	chain.AddTransaction(walletB.BlockchainAddress(), walletA.BlockchainAddress(), 1.0, tx2.GenerateSignature(), walletB.PublicKey())

	total := chain.CalculateTotalAmount(walletA.BlockchainAddress())
	if total != 0 {
		t.Error("Total amount should equal 0")
	}
}
