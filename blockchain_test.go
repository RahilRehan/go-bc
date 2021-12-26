package gobc_test

import (
	"testing"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/stretchr/testify/require"
)

func TestBlockChainCreation(t *testing.T) {
	bc := gobc.NewBlockChain()
	bc.AddBlock("First Block!")
	bc.AddBlock("Second Block!")

	blocks := bc.Blocks

	require.NotNil(t, bc)
	for i, block := range blocks {
		if i == 0 {
			require.NotNil(t, block)
		} else {
			require.Equal(t, block.PrevHash, blocks[i-1].Hash)
		}
	}
}
