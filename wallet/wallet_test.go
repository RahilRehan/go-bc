package wallet_test

import (
	"fmt"
	"testing"

	"github.com/RahilRehan/go-bc/wallet"
)

func TestWalletCreation(t *testing.T) {
	w := wallet.NewWallet()
	fmt.Println(w)
}
