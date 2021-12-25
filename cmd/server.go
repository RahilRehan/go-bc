package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/gorilla/websocket"
)

type server struct {
	port    string
	clients []*client
	txPool  gobc.TransactionPool
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newServer(port string) *server {
	return &server{
		port:   port,
		txPool: *gobc.NewTransactionPool(),
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

	http.HandleFunc("/transaction", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGet(s, rw, r)
		case "POST":
			handlePost(s, rw, r)
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
		fmt.Println(bc)
		s.broadcastMessage(serverLog, c.id)
	}
}

func handlePost(s *server, rw http.ResponseWriter, r *http.Request) {
	var tx gobc.Transaction
	bs, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bs, &tx)
	fmt.Println(tx)
	if err != nil {
		http.Error(rw, "couldn't decode transaction", http.StatusBadRequest)
		return
	}
	s.txPool.Add(&tx)
}

func handleGet(s *server, rw http.ResponseWriter, r *http.Request) {
	bs, _ := json.Marshal(s.txPool.Transactions)
	rw.Write(bs)
}

func (s *server) broadcastMessage(msg, id string) {
	for _, c := range s.clients {
		if c.id != id {
			c.wsConn.WriteMessage(1, []byte(msg))
		}
	}
}
