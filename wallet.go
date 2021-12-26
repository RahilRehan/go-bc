package gobc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
)

const INITIAL_BALANCE = 500

type Wallet struct {
	Balance        int64 `json:"balance"`
	keyPair        *ecdsa.PrivateKey
	PublicKeyECDSA ecdsa.PublicKey
	PublicKey      string `json:"public_key"`
}

func NewWallet() Wallet {

	curve := elliptic.P256()

	keyPair, err := ecdsa.GenerateKey(curve, rand.Reader) // this generates a public & private key pair
	if err != nil {
		log.Fatalln("Error generating key pair: ", err)
	}
	pubkey := keyPair.PublicKey

	wallet := Wallet{
		Balance:        INITIAL_BALANCE,
		keyPair:        keyPair,
		PublicKeyECDSA: pubkey,
		PublicKey:      getPublicKeyString(pubkey),
	}
	return wallet
}

func (w Wallet) String() string {
	return fmt.Sprintf("Wallet: \n Balance: %d\n PublicKey: %v\n", w.Balance, w.PublicKey)
}

func getPublicKeyString(p ecdsa.PublicKey) string {
	return fmt.Sprintf("%x", elliptic.Marshal(p, p.X, p.Y))
}

func (w Wallet) ID() string {
	return fmt.Sprintf("%x", w.keyPair.D.Bytes())
}

func (w Wallet) GetPublicKey() string {
	return fmt.Sprintf("%x", elliptic.Marshal(w.PublicKeyECDSA, w.PublicKeyECDSA.X, w.PublicKeyECDSA.Y))
}

func (w *Wallet) CreateTransaction(recipient *Wallet, amount int64, txp *TransactionPool) *Transaction {
	tx := NewTransaction(w, recipient, amount)
	txp.Add(tx)
	return tx
}
