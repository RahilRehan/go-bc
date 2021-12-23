package blockchain

type blockchain struct {
	blocks []*block
}

func (bc *blockchain) AddBlock(data string) *block {
	prevHash := bc.blocks[len(bc.blocks)-1].GetHash()
	newBlock := MineBlock(prevHash, []byte(data))
	bc.blocks = append(bc.blocks, newBlock)
	return newBlock
}

func (bc *blockchain) GetBlocks() []*block {
	return bc.blocks
}

func NewBlockChain() *blockchain {
	return &blockchain{[]*block{GenesisBlock([]byte("Genesis!"))}}
}
