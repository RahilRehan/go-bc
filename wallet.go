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
	Balance   int64
	keyPair   *ecdsa.PrivateKey
	PublicKey ecdsa.PublicKey
}

func NewWallet() Wallet {

	curve := elliptic.P256()

	keyPair, err := ecdsa.GenerateKey(curve, rand.Reader) // this generates a public & private key pair
	if err != nil {
		log.Fatalln("Error generating key pair: ", err)
	}
	pubkey := keyPair.PublicKey

	wallet := Wallet{
		Balance:   INITIAL_BALANCE,
		keyPair:   keyPair,
		PublicKey: pubkey,
	}
	return wallet
}

func (w Wallet) String() string {
	return fmt.Sprintf("Wallet: \n Balance: %d\n PublicKey: %v\n", w.Balance, w.PublicKey)
}
