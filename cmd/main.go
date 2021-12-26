package main

import (
	"flag"
)

func main() {
	server := flag.Bool("p2pServer", false, "Start in server mode")
	peer := flag.Bool("peer", false, "Start in peer mode")
	serverPort := flag.String("p2pServer-port", ":8080", "Port to listen on")
	flag.Parse()

	if *server {
		s := newServer(*serverPort)
		s.start()
	}

	// create a peer and try to connect to server via /ws
	if *peer {
		p := newPeer(*serverPort)
		p.connectToServer()
	}

}
