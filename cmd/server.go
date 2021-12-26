package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/gorilla/websocket"
)

type server struct {
	port                  string
	clients               []*client
	txPool                *gobc.TransactionPool
	wallets               []gobc.Wallet
	blockchain            gobc.Blockchain
	completedTransactions []gobc.Transaction
}

type txRequest struct {
	Amount      int64  `json:"amount"`
	FromPrivKey string `json:"fromPrivKey"`
	ToPublicKey string `json:"toPbKey"`
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newServer(port string) *server {
	return &server{
		port:   port,
		txPool: gobc.NewTransactionPool(),
	}
}

func (s *server) start() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "couldn't upgrade", http.StatusUpgradeRequired)
			return
		}
		c := newClient(conn)
		s.clients = append(s.clients, c)
		go s.handleWsConn(c)
	})

	http.HandleFunc("/wallet", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/wallet") {
			switch r.Method {
			case "GET":
				s.handleGetWallet(w, r)
			case "POST":
				s.handlePostWallet(w, r)
			}
		}
	})

	http.HandleFunc("/gobc/", func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/gobc/transactions") {
			switch r.Method {
			case "GET":
				// s.handleGetTransaction(rw, r)
			case "POST":
				s.handlePostTransaction(rw, r)
			}
		}
	})

	fmt.Println("=====================================")
	log.Println("P2P Server started on port " + s.port)
	fmt.Println("=====================================")
	log.Fatalln(http.ListenAndServe(s.port, nil))

}

func (s *server) handleWsConn(c *client) {
	defer c.wsConn.Close()
	for {
		_, msg, err := c.wsConn.ReadMessage()
		if err != nil {
			serverLog := fmt.Sprintf("%s# disconnected", c.id)
			log.Println(serverLog)
			s.broadcastMessage(serverLog, c.id)
			break
		}
		// serverLog := fmt.Sprint(string(msg))

		var bc gobc.Blockchain
		err = json.Unmarshal(msg, &bc)
		if err != nil {
			log.Fatalln("Error un marshalling blockchain: ", err)
		}

		s.blockchain = bc
		s.blockchain.Blocks[len(s.blockchain.Blocks)-1].Transactions = s.txPool.Transactions
		s.txPool.Transactions = gobc.NewTransactionPool().Transactions
		bs, err := json.Marshal(s.blockchain)
		if err != nil {
			log.Fatalln("Error marshalling blockchain: ", err)
		}
		fmt.Println(s.txPool.Transactions)
		s.broadcastMessage(string(bs), "everyone")
	}
}

func (s *server) handlePostTransaction(rw http.ResponseWriter, r *http.Request) {
	var txReq txRequest
	bs, _ := io.ReadAll(r.Body)
	json.Unmarshal(bs, &txReq)

	var sender *gobc.Wallet
	var recipient *gobc.Wallet
	for _, w := range s.wallets {
		if w.ID() == txReq.FromPrivKey {
			sender = &w
			break
		}
	}

	for _, w := range s.wallets {
		if w.GetPublicKey() == txReq.ToPublicKey {
			recipient = &w
			break
		}
	}

	if sender == nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Sender wallet not found"))
		return
	}

	if recipient == nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Recipient wallet not found"))
		return
	}

	if sender.ID() == recipient.ID() {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("You cannot send to yourself"))
		return
	}

	tx := gobc.NewTransaction(sender, recipient, txReq.Amount)
	tx.Sign(sender, &tx.Output)
	s.txPool.Add(tx)

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("Transaction added to unverified and incomplete transactions pool! Will be confirmed in next available block."))
}

func (s *server) handleGetWallet(w http.ResponseWriter, r *http.Request) {
	walletPubKeys := make(map[int]string, 0)
	i := 0
	for _, w := range s.wallets {
		walletPubKeys[i] = w.GetPublicKey()
		i++
	}
	bs, _ := json.Marshal(walletPubKeys)
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func (s *server) handlePostWallet(w http.ResponseWriter, r *http.Request) {
	wallet := gobc.NewWallet()
	s.wallets = append(s.wallets, wallet)
	wallet.ID()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Wallet is created \nPlease remember this secret id for making transactions \nSecret ID: %s", wallet.ID())))
}

func (s *server) broadcastMessage(msg, id string) {
	for _, c := range s.clients {
		if c.id != id {
			c.wsConn.WriteMessage(1, []byte(msg))
		}
	}
}
