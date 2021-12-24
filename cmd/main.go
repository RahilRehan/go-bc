package main

import (
	"flag"
)

func main() {
	server := flag.Bool("server", false, "Start in server mode")
	peer := flag.Bool("peer", false, "Start in peer mode")
	port := flag.String("port", ":8080", "Port to listen on")
	flag.Parse()

	if *server {
		s := newServer(*port)
		s.start()
	}

	// create a peer and try to connect to server via /ws
	if *peer {
		p := newPeer(*port)
		p.connectToServer()
	}

}
