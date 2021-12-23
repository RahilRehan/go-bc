package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Notes:
// do we need an interface?
// getters and setters?
// I don't like the idea of getters as everything in blockchain is public.
// Doubt: if I make my struct fields exposed or expose a setter, anyone can change my fields (data can be tampered) - Integrity?
// Solution: Just have a constructor via which block can be created. And have getters for each fields but no setters as we don't want to allow data tampering.

const HASH_SIZE = 32

type block struct {
	timestamp time.Time
	prevHash  [HASH_SIZE]byte
	hash      [HASH_SIZE]byte
	data      []byte
}

// Stringer

func (b *block) String() string {
	return fmt.Sprintf("Timestamp: %s\nPrevHash: %x\nHash: %x\nData: %s\n",
		b.timestamp.String(), string(b.prevHash[:]), string(b.hash[:]), string(b.data))
}

// Getters

func (b *block) GetTimestamp() time.Time {
	return b.timestamp
}

func (b *block) GetPrevHash() [HASH_SIZE]byte {
	return b.prevHash
}

func (b *block) GetHash() [HASH_SIZE]byte {
	return b.hash
}

func (b *block) GetData() []byte {
	return b.data
}

// Constructors

// Create a Genesis Block initially
func GenesisBlock(data []byte) *block {
	return MineBlock(sha256.Sum256([]byte("---------")), data)
}

// Mine a new block based out of the previous block
func MineBlock(prevHash [HASH_SIZE]byte, data []byte) *block {

	block := block{
		timestamp: time.Now(),
		prevHash:  prevHash,
		data:      []byte(data),
	}

	block.hash = sha256.Sum256([]byte(block.timestamp.String() + string(block.prevHash[:]) + string(block.data)))

	return &block
}
