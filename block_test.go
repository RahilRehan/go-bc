package gobc_test

import (
	"testing"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/stretchr/testify/require"
)

func TestGenesisBlock(t *testing.T) {
	block := gobc.GenesisBlock([]byte("Genesis!"))

	require.NotNil(t, block)
}

func TestMineBlock(t *testing.T) {

	genisisBlock := gobc.GenesisBlock([]byte("Genesis!"))
	prevHash := genisisBlock.GetHash()
	block := gobc.MineBlock(prevHash, []byte("First Block!"))

	require.NotNil(t, block)
	require.Equal(t, block.GetPrevHash(), prevHash)
	require.Equal(t, block.GetData(), []byte("First Block!"))
}
