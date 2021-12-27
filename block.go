package gobc

import (
	"bytes"
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
const MINE_RATE = (2 * time.Second)

var DIFFICULTY = 4

type Block struct {
	Timestamp    time.Time     `json:"timestamp"`
	PrevHash     string        `json:"prevHash"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
	Nonce        int64         `json:"nonce"`
}

func (b Block) String() string {
	return fmt.Sprintf(`
		Timestamp: %v
		PrevHash: %s
		Hash: %s
		Transactions: %v
		Nonce: %d
		
	`, b.Timestamp, b.PrevHash, b.Hash, b.Transactions, b.Nonce)
}

// Create a Genesis Block initially
func GenesisBlock() Block {
	dummyHash := NewSHA256([]byte("---------"))
	return MineBlock(string(dummyHash[:]), 0, time.Now())
}

// Mine a new block based out of the previous block
func MineBlock(prevHash string, nonce int64, prevBlockCreatedTime time.Time) Block {

	var block Block
	var hash string

	for {
		block = Block{
			Timestamp: time.Now(),
			PrevHash:  string(prevHash[:]),
			Nonce:     nonce,
		}

		hash = NewSHA256([]byte(block.Timestamp.String() + block.PrevHash + fmt.Sprint(nonce)))
		if hash[:DIFFICULTY] == string(bytes.Repeat([]byte("0"), DIFFICULTY)) {
			if time.Since(prevBlockCreatedTime) < MINE_RATE {
				DIFFICULTY += 1
				if DIFFICULTY > 64 {
					DIFFICULTY = 64
				}
			} else if time.Since(prevBlockCreatedTime) > MINE_RATE {
				DIFFICULTY -= 1
				if DIFFICULTY < 1 {
					DIFFICULTY = 1
				}
			}
			break
		}

		nonce++
	}

	block.Hash = string(hash[:])

	return block
}
