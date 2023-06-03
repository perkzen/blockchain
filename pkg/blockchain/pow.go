package blockchain

import (
	"fmt"
	"strings"
)

// MINING_DIFFICULTY
// determines the number of zeros on the beginning of the hash
const (
	MINING_DIFFICULTY = 4
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 0.1
)

type PoW struct {
	Block *Block
}

func NewProofOfWork(b *Block) *PoW {
	return &PoW{
		Block: b,
	}
}

func (p *PoW) calculateHash(nonce int) [32]byte {
	p.Block.Nonce = nonce
	return p.Block.Hash()
}

func (p *PoW) Proof() (nonce int) {
	nonce = 0
	for !p.Validate(nonce) {
		nonce++
	}
	return nonce
}

func isValid(hash [32]byte) bool {
	guess := fmt.Sprintf("%x", hash)
	zeros := strings.Repeat("0", MINING_DIFFICULTY)
	return guess[:MINING_DIFFICULTY] == zeros
}

func (p *PoW) Validate(nonce int) bool {
	hash := p.calculateHash(nonce)
	return isValid(hash)
}
