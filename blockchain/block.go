package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
const DIFFICULTY = 5

type Block struct {
	Timestamp time.Time `json:"timestamp"`
	PrevHash  string    `json:"prevHash"`
	Hash      string    `json:"hash"`
	Data      []byte    `json:"data"`
	Nonce     int       `json:"nonce"`
}

// Constructors

// Create a Genesis Block initially
func GenesisBlock(data []byte) Block {
	dummyHash := NewSHA256([]byte("---------"))
	return MineBlock(string(dummyHash[:]), data)
}

// Mine a new block based out of the previous block
func MineBlock(prevHash string, data []byte) Block {

	var block Block
	var hash string
	nonce := 0

	for {
		block = Block{
			Timestamp: time.Now(),
			PrevHash:  string(prevHash[:]),
			Data:      []byte(data),
			Nonce:     nonce,
		}

		hash = NewSHA256([]byte(block.Timestamp.String() + block.PrevHash + string(block.Data) + fmt.Sprint(nonce)))
		if hash[:DIFFICULTY] == string(bytes.Repeat([]byte("0"), DIFFICULTY)) {
			break
		}
		nonce++
	}

	block.Hash = string(hash[:])

	return block
}

func NewSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
