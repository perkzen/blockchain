package blockchain

import (
	"strings"
)

type PoW struct {
	Block *Block
}

func NewProofOfWork(b *Block) *PoW {
	return &PoW{
		Block: b,
	}
}

func (p *PoW) calculateHash(nonce int) string {
	p.Block.Nonce = nonce
	return p.Block.Hash()
}

func (p *PoW) Validate() (nonce int) {
	nonce = 0
	for !p.IsValid(nonce) {
		nonce++
	}
	return nonce
}

func (p *PoW) IsValid(nonce int) bool {
	hash := p.calculateHash(nonce)
	zeros := strings.Repeat("0", MINING_DIFFICULTY)
	return hash[:MINING_DIFFICULTY] == zeros
}
