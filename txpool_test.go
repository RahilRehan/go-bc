package gobc_test

import (
	"testing"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/stretchr/testify/require"
)

func TestTransactionPoolCreation(t *testing.T) {
	txp := gobc.NewTransactionPool()
	require.NotNil(t, txp)
}

func TestTransactionPoolAddNewTransaction(t *testing.T) {
	sender := gobc.NewWallet()
	recipient := gobc.NewWallet()
	amount := int64(70)

	txp := gobc.NewTransactionPool()
	prevTxpLen := len(txp.Transactions)
	tx := gobc.NewTransaction(&sender, &recipient, amount)
	tx.Sign(&sender, &tx.Output)
	txp.Add(tx)

	require.Equal(t, len(txp.Transactions), prevTxpLen+1)
}
