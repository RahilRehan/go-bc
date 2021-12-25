package gobc_test

import (
	"testing"
	"time"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/stretchr/testify/require"
)

func TestGenesisBlock(t *testing.T) {
	block := gobc.GenesisBlock([]byte("Genesis!"))

	require.NotNil(t, block)
}

func TestMineBlock(t *testing.T) {

	genisisBlock := gobc.GenesisBlock([]byte("Genesis!"))
	prevHash := genisisBlock.Hash
	block := gobc.MineBlock(prevHash, []byte("First Block!"), 0, time.Now())

	require.NotNil(t, block)
	require.Equal(t, block.PrevHash, prevHash)
	require.Equal(t, block.Data, []byte("First Block!"))
}
