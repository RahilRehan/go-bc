package gobc_test

import (
	"fmt"
	"testing"

	gobc "github.com/RahilRehan/go-bc"
)

func TestWalletCreation(t *testing.T) {
	w := gobc.NewWallet()
	fmt.Println(w)
}
