package blockchain_test

import (
	"testing"

	"github.com/RahilRehan/go-bc/blockchain"
	"github.com/stretchr/testify/require"
)

func TestGenesisBlock(t *testing.T) {
	block := blockchain.GenesisBlock([]byte("Genesis!"))

	require.NotNil(t, block)
}

func TestMineBlock(t *testing.T) {

	genisisBlock := blockchain.GenesisBlock([]byte("Genesis!"))
	prevHash := genisisBlock.GetHash()
	block := blockchain.MineBlock(prevHash, []byte("First Block!"))

	require.NotNil(t, block)
	require.Equal(t, block.GetPrevHash(), prevHash)
	require.Equal(t, block.GetData(), []byte("First Block!"))
}
