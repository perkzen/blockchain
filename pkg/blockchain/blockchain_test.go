package blockchain

import (
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
	chain := InitBlockchain()
	if len(chain.Blocks) < 1 || len(chain.Blocks) > 1 {
		t.Error("Chain should have only 1 block on initiation")
	}
}

func TestBlockchain_AddBlock(t *testing.T) {
	chain := InitBlockchain()
	chain.AddBlock()
	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
}

func TestBlockchain_AddTransaction(t *testing.T) {
	chain := InitBlockchain()
	chain.AddTransaction("A", "B", 1)
	if len(chain.TxPool) < 1 {
		t.Error("Chain should have 1 transaction in pool")
	}
}

func TestBlockchain_ClearPool(t *testing.T) {
	chain := InitBlockchain()
	chain.AddTransaction("A", "B", 1)
	chain.AddBlock()
	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
	if len(chain.TxPool) >= 1 {
		t.Error("Transaction pool should be empty")
	}
}
