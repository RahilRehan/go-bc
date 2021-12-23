package main

import (
	"github.com/RahilRehan/go-bc/blockchain"
	"github.com/RahilRehan/go-bc/p2p"
)

func main() {

	// Create a blockchain and assign to p2p blockchain
	bc := blockchain.NewBlockChain()
	p2p.MyBlockchain = bc

	p2p.MyBlockchain = blockchain.NewBlockChain()
	host := *p2p.CreateHost("9595")
	host.SetStreamHandler("/blockchain", p2p.StreamHandler)

	select {}
}
