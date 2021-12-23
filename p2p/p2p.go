package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RahilRehan/go-bc/blockchain"
	libp2p "github.com/libp2p/go-libp2p"
	net "github.com/libp2p/go-libp2p-core"
	host "github.com/libp2p/go-libp2p-core/host"
	ma "github.com/multiformats/go-multiaddr"
)

var mutex = &sync.Mutex{}
var MyBlockchain *blockchain.Blockchain

func CreateHost(port string) *host.Host {
	host, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", port)))
	if err != nil {
		log.Fatalln("Error creating the host: ", err)
	}

	multiAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", host.ID().Pretty()))
	if err != nil {
		log.Fatalln("Error creating the multiaddress: ", err)
	}

	addr := host.Addrs()[0]
	fullAddr := addr.Encapsulate(multiAddr)

	log.Printf("I am %s\n", fullAddr)
	log.Printf("Run \"go run p2p.go -l %s -d %s\" on a different terminal\n", port+fmt.Sprint(1), fullAddr)
	return &host
}

func StreamHandler(s net.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readFromPeers(rw)
	go writeToPeers(rw)
}

func readFromPeers(rw *bufio.ReadWriter) {
	// infinite loop, keep reading from the stream
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatalln("Error reading data: ", err)
		}
		if str == "" {
			return
		}
		if str != "\n" {
			chain := make([]*blockchain.Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()

			if len(chain) > len(MyBlockchain.Blocks) {
				MyBlockchain.Blocks = chain
				log.Println("Received blockchain length: ", len(chain))
				log.Panicln("Updated the local blockchain")
			}

			mutex.Unlock()
		}
	}
}

func writeToPeers(rw *bufio.ReadWriter) {
	// separate infinite go routine to write current blockchain to the stream every 5 seconds.
	go func() {
		for {
			mutex.Lock()
			bs, err := json.Marshal(MyBlockchain.Blocks)
			if err != nil {
				log.Fatalln("Error marshalling: ", err)
			}
			mutex.Unlock()

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bs)))
			rw.Flush()
			mutex.Unlock()
			time.Sleep(5 * time.Second)
		}
	}()

}
