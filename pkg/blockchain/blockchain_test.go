package blockchain

import (
	"bytes"
	"testing"
)

func TestCreateGenesisBlock(t *testing.T) {
	genesis := CreateGenesisBlock()
	if !bytes.Equal(genesis.PrevHash, []byte{}) {
		t.Error("Previous hash should be empty")
	}
}

func TestInitBlockchain(t *testing.T) {
	chain := InitBlockchain()
	if len(chain.Blocks) < 1 || len(chain.Blocks) > 1 {
		t.Error("Chain should have only 1 block on initiation")
	}
	if !bytes.Equal(chain.Blocks[0].PrevHash, []byte{}) {
		t.Error("Previous hash should be empty because it's a GENESIS block")
	}
}
