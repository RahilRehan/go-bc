package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/gorilla/websocket"
)

type server struct {
	port    string
	clients []*client
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newServer(port string) *server {
	return &server{
		port: port,
	}
}

func (s *server) start() {

	fmt.Println("=====================================")
	log.Println("Server started on port " + s.port)
	fmt.Println("=====================================")

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

	http.ListenAndServe(s.port, nil)
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
			log.Fatalln("Error unmarshalling blockchain: ", err)
		}
		fmt.Println(bc)
		s.broadcastMessage(serverLog, c.id)
	}
}

func (s *server) broadcastMessage(msg, id string) {
	for _, c := range s.clients {
		if c.id != id {
			c.wsConn.WriteMessage(1, []byte(msg))
		}
	}
}
