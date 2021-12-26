package gobc_test

import (
	"testing"
	"time"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/stretchr/testify/require"
)

func TestGenesisBlock(t *testing.T) {
	block := gobc.GenesisBlock()

	require.NotNil(t, block)
}

func TestMineBlock(t *testing.T) {

	genisisBlock := gobc.GenesisBlock()
	prevHash := genisisBlock.Hash
	block := gobc.MineBlock(prevHash, 0, time.Now())

	require.NotNil(t, block)
	require.Equal(t, block.PrevHash, prevHash)
}
