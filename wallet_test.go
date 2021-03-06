package gobc_test

import (
	"fmt"
	"testing"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/stretchr/testify/require"
)

func TestWalletCreation(t *testing.T) {
	w := gobc.NewWallet()
	fmt.Println(w)
}

func TestTransactionCreationOnWalletAndAddTransactionToTransactionPool(t *testing.T) {
	sender := gobc.NewWallet()
	recipient := gobc.NewWallet()
	amount := int64(70)

	txp := gobc.NewTransactionPool()
	prevTxpLen := len(txp.Transactions)
	tx := sender.CreateTransaction(&recipient, amount, txp)
	tx.Sign(&sender, &tx.Output)
	require.Equal(t, prevTxpLen+1, len(txp.Transactions))
}
