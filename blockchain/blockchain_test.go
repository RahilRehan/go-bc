package blockchain_test

import (
	"testing"

	"github.com/RahilRehan/go-bc/blockchain"
	"github.com/stretchr/testify/require"
)

func TestBlockChainCreation(t *testing.T) {
	bc := blockchain.NewBlockChain()
	bc.AddBlock("First Block!")
	bc.AddBlock("Second Block!")

	blocks := bc.GetBlocks()

	require.NotNil(t, bc)
	for i, block := range blocks {
		if i == 0 {
			require.Equal(t, block.GetData(), []byte("Genesis!"))
		} else {
			require.Equal(t, block.GetPrevHash(), blocks[i-1].GetHash())
		}
	}
}
