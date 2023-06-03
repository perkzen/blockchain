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
	chain := InitBlockchain("")
	if len(chain.Blocks) < 1 || len(chain.Blocks) > 1 {
		t.Error("Chain should have only 1 block on initiation")
	}
}

func TestBlockchain_AddBlock(t *testing.T) {
	chain := InitBlockchain("")
	chain.AddBlock()
	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
}

func TestBlockchain_AddTransaction(t *testing.T) {
	chain := InitBlockchain("")
	chain.AddTransaction("A", "B", 1)
	if len(chain.TxPool) < 1 {
		t.Error("Chain should have 1 transaction in pool")
	}
}

func TestBlockchain_ClearPool(t *testing.T) {
	chain := InitBlockchain("")
	chain.AddTransaction("A", "B", 1)
	chain.AddBlock()
	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
	if len(chain.TxPool) >= 1 {
		t.Error("Transaction pool should be empty")
	}
}

func TestBlockchain_Mining(t *testing.T) {
	chain := InitBlockchain("")
	chain.Mining()
	lastBlock := chain.lastBlock()
	if len(chain.Blocks) <= 1 {
		t.Error("Chain should have more than 1 block")
	}
	if len(chain.TxPool) >= 1 {
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
	chain := InitBlockchain("")
	chain.AddTransaction("A", "B", 1)
	chain.AddTransaction("B", "A", 1)
	total := chain.CalculateTotalAmount("A")
	if total != 0 {
		t.Error("Total amount should equal 0")
	}
}
