package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
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

func (p *PoW) getData(nonce int) []byte {
	return bytes.Join([][]byte{
		p.Block.PrevHash,
		p.Block.Data,
		toHex(int64(nonce)),
		toHex(DIFFICULTY),
	}, []byte{})
}

func (p *PoW) Mine() (hash [32]byte, nonce int) {
	nonce = 0
	for {
		data := p.getData(nonce)
		hash := sha256.Sum256(data)

		if isValid(hash) {
			break
		}
		nonce++
	}

	return hash, nonce
}

func isValid(hash [32]byte) bool {
	guess := fmt.Sprintf("%x", hash)
	zeros := strings.Repeat("0", DIFFICULTY)
	return guess[:DIFFICULTY] == zeros
}

func (p *PoW) Validate(nonce int) bool {
	data := p.getData(nonce)
	hash := sha256.Sum256(data)
	return isValid(hash)
}

func toHex(n int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, n)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}
