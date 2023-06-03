package blockchain

import (
	"time"
)

type Block struct {
	Hash      []byte
	PrevHash  []byte
	Data      []byte
	Timestamp int64
	Nonce     int
}

type Blockchain struct {
	Blocks []*Block
}

func createBlock(prevHash []byte, data string) *Block {
	b := &Block{
		Data:      []byte(data),
		PrevHash:  prevHash,
		Timestamp: time.Now().UnixNano(),
		Nonce:     0,
	}

	pow := NewProofOfWork(b)
	hash, nonce := pow.Mine()

	b.Nonce = nonce
	b.Hash = hash[:]

	return b
}

func CreateGenesisBlock() *Block {
	return createBlock([]byte{}, "GENESIS")
}

func (chain *Blockchain) AddBlock(data string) {
	prevHash := chain.lastBlock().PrevHash
	chain.Blocks = append(chain.Blocks, createBlock(prevHash, data))
}

func (chain *Blockchain) lastBlock() *Block {
	return chain.Blocks[len(chain.Blocks)-1]
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{CreateGenesisBlock()}}
}
