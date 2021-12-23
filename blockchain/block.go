package blockchain

import (
	"crypto/sha256"
	"time"
)

// Notes:
// do we need an interface?
// getters and setters?
// I don't like the idea of getters as everything in blockchain is public.
// Doubt: if I make my struct fields exposed or expose a setter, anyone can change my fields (data can be tampered) - Integrity?
// Solution: Just have a constructor via which block can be created. And have getters for each fields but no setters as we don't want to allow data tampering.
// But for un-marshaling Json we need to have exported fields :(

const HASH_SIZE = 32

type block struct {
	Timestamp time.Time       `json:"timestamp"`
	PrevHash  [HASH_SIZE]byte `json:"prevHash"`
	Hash      [HASH_SIZE]byte `json:"hash"`
	Data      []byte          `json:"data"`
}

// Constructors

// Create a Genesis Block initially
func GenesisBlock(data []byte) *block {
	return MineBlock(sha256.Sum256([]byte("---------")), data)
}

// Mine a new block based out of the previous block
func MineBlock(prevHash [HASH_SIZE]byte, data []byte) *block {

	block := block{
		Timestamp: time.Now(),
		PrevHash:  prevHash,
		Data:      []byte(data),
	}

	block.Hash = sha256.Sum256([]byte(block.Timestamp.String() + string(block.PrevHash[:]) + string(block.Data)))

	return &block
}
