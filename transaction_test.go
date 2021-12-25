package gobc_test

import (
	"testing"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/stretchr/testify/require"
)

func TestTransactionCreation(t *testing.T) {
	sender := gobc.NewWallet()
	require.NotNil(t, sender)
}

func TestNegativeAmountTransaction(t *testing.T) {
	sender := gobc.NewWallet()
	recipient := gobc.NewWallet()
	amount := int64(-20)
	tx := gobc.NewTransaction(&sender, recipient.PublicKey, amount)

	require.Nil(t, tx)
}

func TestInSufficientBalanceTransaction(t *testing.T) {
	sender := gobc.NewWallet()
	recipient := gobc.NewWallet()
	amount := int64(sender.Balance + 1)
	tx := gobc.NewTransaction(&sender, recipient.PublicKey, amount)

	require.Nil(t, tx)
}

func TestTransactionBetweenTwoWallets(t *testing.T) {
	sender := gobc.NewWallet()
	recipient := gobc.NewWallet()
	amount := int64(20)
	tx := gobc.NewTransaction(&sender, recipient.PublicKey, amount)
	tx.Sign(&sender, &tx.Output)

	require.NotNil(t, tx)
	require.Equal(t, tx.Output[0].Address, sender.PublicKey)
	require.Equal(t, tx.Output[0].Amount, sender.Balance-amount)

	require.Equal(t, tx.Output[1].Address, recipient.PublicKey)
	require.Equal(t, tx.Output[1].Amount, amount)

	verified := tx.Verify(&sender, tx.Input.Signature)
	require.True(t, verified)
}
