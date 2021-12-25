package blockchain_test

import (
	"testing"
	"time"

	"github.com/RahilRehan/go-bc/blockchain"
	"github.com/stretchr/testify/require"
)

func TestGenesisBlock(t *testing.T) {
	block := blockchain.GenesisBlock([]byte("Genesis!"))

	require.NotNil(t, block)
}

func TestMineBlock(t *testing.T) {

	genisisBlock := blockchain.GenesisBlock([]byte("Genesis!"))
	prevHash := genisisBlock.Hash
	block := blockchain.MineBlock(prevHash, []byte("First Block!"), 0, time.Now())

	require.NotNil(t, block)
	require.Equal(t, block.PrevHash, prevHash)
	require.Equal(t, block.Data, []byte("First Block!"))
}
