package gobc

import "fmt"

type Blockchain struct {
	Blocks []Block
}

// Add a new block to the blockchain
func (bc Blockchain) AddBlock(data string) Blockchain {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	prevHash := prevBlock.Hash
	newBlock := MineBlock(prevHash, []byte(data), 0, prevBlock.Timestamp)
	bc.Blocks = append(bc.Blocks, newBlock)
	return bc
}

func (bc Blockchain) String() string {
	blocks := ""
	for _, block := range bc.Blocks {
		blocks += fmt.Sprintln(block)
	}
	return fmt.Sprintln("================BLOCKCHAIN================") +
		fmt.Sprint(blocks) +
		fmt.Sprintln("==========================================")
}

// Return new blockchain with genesis block
func NewBlockChain() Blockchain {
	return Blockchain{
		[]Block{
			GenesisBlock([]byte("Genesis!")),
		},
	}
}
