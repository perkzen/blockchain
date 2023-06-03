package blockchain

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Hash      []byte
	PrevHash  []byte
	Data      []byte
	Timestamp int64
}

type Blockchain struct {
	Blocks []*Block
}

func (b *Block) generateBlockHash() {
	data := bytes.Join([][]byte{b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}

func createBlock(prevHash []byte, data string) *Block {
	b := &Block{
		Data:      []byte(data),
		PrevHash:  prevHash,
		Timestamp: time.Now().Unix(),
	}

	b.generateBlockHash()

	return b
}

func CreateGenesisBlock() *Block {
	return createBlock([]byte{}, "GENESIS")
}

func (chain *Blockchain) AddBlock(data string) {
	prevHash := chain.Blocks[len(chain.Blocks)-1].PrevHash
	chain.Blocks = append(chain.Blocks, createBlock(prevHash, data))
}

func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{CreateGenesisBlock()}}
}
