package blockchain

import (
	"fmt"
	"strings"
)

// DIFFICULTY
// determines the number of zeros on the beginning of the hash
const DIFFICULTY = 4

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

func (p *PoW) Mine() (nonce int) {
	nonce = 0
	for {
		hash := p.calculateHash(nonce)
		if isValid(hash) {
			break
		}
		nonce++
	}
	return nonce
}

func isValid(hash [32]byte) bool {
	guess := fmt.Sprintf("%x", hash)
	zeros := strings.Repeat("0", DIFFICULTY)
	return guess[:DIFFICULTY] == zeros
}

func (p *PoW) Validate(nonce int) bool {
	hash := p.calculateHash(nonce)
	return isValid(hash)
}
