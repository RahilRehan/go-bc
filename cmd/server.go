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
	port       string
	clients    []*client
	txPool     *gobc.TransactionPool
	wallets    []gobc.Wallet
	blockchain gobc.Blockchain
}

type txRequest struct {
	Amount int64 `json:"amount"`
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

	http.HandleFunc("/gobc/", func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/gobc/transactions") {
			switch r.Method {
			case "GET":
				s.handleGet(rw, r)
			case "POST":
				s.handlePost(rw, r)
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
		serverLog := fmt.Sprint(string(msg))

		var bc gobc.Blockchain
		err = json.Unmarshal(msg, &bc)
		if err != nil {
			log.Fatalln("Error un marshalling blockchain: ", err)
		}
		s.blockchain = bc
		s.broadcastMessage(serverLog, c.id)
	}
}

func (s *server) handlePost(rw http.ResponseWriter, r *http.Request) {
	sender := gobc.NewWallet()
	receiver := gobc.NewWallet()
	s.wallets = append(s.wallets, sender)
	s.wallets = append(s.wallets, receiver)

	var txReq txRequest
	bs, _ := io.ReadAll(r.Body)
	json.Unmarshal(bs, &txReq)
	tx := gobc.NewTransaction(&sender, &receiver, txReq.Amount)
	s.txPool.Add(tx)

	fmt.Println("TX POOL: ", s.txPool)

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("Transaction added to unverified transactions pool! Will be confirmed in next available block."))

}

func (s *server) handleGet(rw http.ResponseWriter, r *http.Request) {
	// for _, block := range s.blockchain.Blocks {
	// 	for _, tx := range block.Transactions {
	// 		fmt.Println(tx)
	// 	}
	// }
	rw.Write([]byte("These are all the validated transactions!"))
}

func (s *server) broadcastMessage(msg, id string) {
	for _, c := range s.clients {
		if c.id != id {
			c.wsConn.WriteMessage(1, []byte(msg))
		}
	}
}
