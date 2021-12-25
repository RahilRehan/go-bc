package main

import (
	"flag"

	"github.com/RahilRehan/go-bc/web"
)

func main() {
	webServer := flag.Bool("webserver", false, "Run the web client")
	server := flag.Bool("p2pServer", false, "Start in server mode")
	peer := flag.Bool("peer", false, "Start in peer mode")
	webServerPort := flag.String("webserver-port", ":9090", "Port to listen on")
	serverPort := flag.String("p2pServer-port", ":8080", "Port to listen on")
	flag.Parse()

	if *webServer {
		web.NewWebServer(*webServerPort)
	}

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
