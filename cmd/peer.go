package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	gobc "github.com/RahilRehan/go-bc"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// each peer has a blockchain attached to it, it connects to server and handles broadcasts
type peer struct {
	host       string
	blockchain gobc.Blockchain
}

func newPeer(port string) *peer {
	return &peer{
		host:       "localhost" + port,
		blockchain: gobc.NewBlockChain(),
	}
}

// sends a request to server to join into p2p connection
func (p *peer) connectToServer() {
	wsServer := url.URL{Scheme: "ws", Host: p.host, Path: "/ws"}
	conn, res, err := websocket.DefaultDialer.Dial(wsServer.String(), nil)
	if err != nil {
		log.Fatalf("couldn't connect to websocket server: %s\n", err.Error())
	}
	if res.StatusCode == http.StatusUpgradeRequired {
		log.Fatalf("the server did not upgrade to ws protocol\n")
	}

	go p.handleNewMessage(conn)

	for {
		bs, err := json.Marshal(p.blockchain)
		if err != nil {
			log.Fatalln("couldn't marshal blockchain")
		}
		err = conn.WriteMessage(1, bs)
		if err != nil {
			log.Fatalln("couldn't send message: ", err.Error())
		}

		p.blockchain = p.blockchain.AddBlock("new data")
	}

	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	msg := scanner.Text()
	// 	if len(msg) > 0 {
	// 		err := conn.WriteMessage(1, []byte(scanner.Text()))
	// 		if err != nil {
	// 			fmt.Printf("couldn't send message: %s\n", err.Error())
	// 		}
	// 	}
	// 	fmt.Print("~> ")
	// }
}

//handles broadcasted messages
func (p *peer) handleNewMessage(conn *websocket.Conn) {

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("error while reading: %s\n", err.Error())
		}
		var receivedBlockchain gobc.Blockchain
		err = json.Unmarshal(msg, &receivedBlockchain)
		if err != nil {
			log.Fatalf("error while unmarshaling: %s\n", err.Error())
		}

		// sync blockchain according to received blockchains length, longest blockchain wins
		if len(receivedBlockchain.Blocks) > len(p.blockchain.Blocks) {
			p.blockchain = receivedBlockchain
			log.Println("=========> Blockchain updated!!!")
		}

	}
}

// when server adds peer to its pool, server creates a client to interact with peer. It has websocket connection object
type client struct {
	id     string
	wsConn *websocket.Conn
}

// whenever there is a request to /ws on server, we create a new client
func newClient(conn *websocket.Conn) *client {
	return &client{
		wsConn: conn,
		id:     uuid.New().String(),
	}
}
