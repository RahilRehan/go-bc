package blockchain

type Blockchain struct {
	Blocks []Block
}

// Add a new block to the blockchain
func (bc Blockchain) AddBlock(data string) Blockchain {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	newBlock := MineBlock(prevHash, []byte(data))
	bc.Blocks = append(bc.Blocks, newBlock)
	return bc
}

// Return new blockchain with genesis block
func NewBlockChain() Blockchain {
	return Blockchain{[]Block{GenesisBlock([]byte("Genesis!"))}}
}
