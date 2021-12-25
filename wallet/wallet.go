package wallet

const INITIAL_BALANCE = 500

type Wallet struct {
	balance   int64
	keyPair   [2]string
	publicKey string
}

func NewWallet() *Wallet {
	wallet := Wallet{
		balance:   INITIAL_BALANCE,
		keyPair:   [2]string{},
		publicKey: "",
	}
	return &wallet
}
