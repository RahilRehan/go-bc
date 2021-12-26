package gobc

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id       uuid.UUID       `json:"id"`
	Input    InputDetail     `json:"input"`  // contains => details about sender, senders [original balance, signature and public key]
	Output   [2]OutputDetail `json:"output"` // contains => [OutputDetail of senders balance after transaction + senders public key] and [OutputDetail of amount amount received by receiver + recipients public key]
	Verified bool            `json:"verified"`
	Complete bool            `json:"complete"`
}

type OutputDetail struct {
	Amount  int64
	Address ecdsa.PublicKey
}

type InputDetail struct {
	Timestamp time.Time
	Amount    int64
	Address   ecdsa.PublicKey
	Signature SignatureDetail
}

type SignatureDetail struct {
	SignHash []byte
	R, S     *big.Int
}

func (o OutputDetail) String() string {
	return fmt.Sprintf("OutputDetail: \n Amount: %d\n Address: %v\n", o.Amount, o.Address)
}

func NewTransaction(sender *Wallet, recipient *Wallet, amount int64) *Transaction {
	if amount < 0 {
		log.Println("Amount cannot be negative")
		return nil
	}
	if amount > sender.Balance {
		log.Println("Insufficient balance")
		return nil
	}

	output := [2]OutputDetail{
		{
			Amount:  sender.Balance - amount,
			Address: sender.PublicKey,
		},
		{
			Amount:  amount,
			Address: recipient.PublicKey,
		},
	}

	input := InputDetail{
		Timestamp: time.Now(),
		Amount:    sender.Balance,
		Address:   sender.PublicKey,
	}

	return &Transaction{
		Id:       uuid.New(),
		Input:    input,
		Output:   output,
		Verified: false,
		Complete: false,
	}
}

func (tx *Transaction) Sign(senderWallet *Wallet, outputDetail *[2]OutputDetail) {

	h := md5.New()

	io.WriteString(h, outputDetail[0].String()+outputDetail[1].String())
	signhash := h.Sum(nil)

	// sign with senders private key
	r, s, err := ecdsa.Sign(rand.Reader, senderWallet.keyPair, signhash)
	if err != nil {
		log.Fatalln("Error signing: ", err)
	}
	tx.Input.Signature = SignatureDetail{signhash, r, s}
}

// verify the transaction with the public key of the sender
func (tx *Transaction) Verify(senderWallet *Wallet, sd SignatureDetail) bool {
	verifystatus := ecdsa.Verify(&senderWallet.PublicKey, sd.SignHash, sd.R, sd.S)
	return verifystatus
}
