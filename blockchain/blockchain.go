package blockchain

type blockchain struct {
	Blocks []*block
}

// Add a new block to the blockchain
func (bc *blockchain) AddBlock(data string) *block {
	prevHash := bc.Blocks[len(bc.Blocks)-1].Hash
	newBlock := MineBlock(prevHash, []byte(data))
	bc.Blocks = append(bc.Blocks, newBlock)
	return newBlock
}

// Return new blockchain with genesis block
func NewBlockChain() *blockchain {
	return &blockchain{[]*block{GenesisBlock([]byte("Genesis!"))}}
}
