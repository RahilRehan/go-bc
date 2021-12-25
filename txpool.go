package gobc

import (
	"log"
)

const TP_SIZE = 10

type TransactionPool struct {
	Transactions []*Transaction
}

func (txp *TransactionPool) Add(tx *Transaction) {
	if len(txp.Transactions) >= TP_SIZE {
		log.Fatalln("Transaction pool is full")
	}
	txp.Transactions = append(txp.Transactions, tx)
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make([]*Transaction, 1),
	}
}
