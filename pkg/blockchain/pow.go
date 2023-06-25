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

func (p *PoW) Run() (nonce int) {
	nonce = 0
	for !p.IsValid(nonce) {
		nonce++
	}
	return nonce
}

func (p *PoW) IsValid(nonce int) bool {
	p.Block.setNonce(nonce)
	hash := p.Block.Hash()
	zeros := strings.Repeat("0", MINING_DIFFICULTY)
	return hash[:MINING_DIFFICULTY] == zeros
}
