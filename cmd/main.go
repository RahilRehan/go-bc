package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RahilRehan/go-bc/blockchain"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to GO Blockchain!"))
	})

	// Create a blockchain
	blockchain := blockchain.NewBlockChain()

	http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			return
		} else if r.Method == "GET" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if len(blockchain.Blocks) == 0 {
				w.Write([]byte("No blocks in the chain"))
			} else {
				bs, err := json.Marshal(blockchain.Blocks)
				if err != nil {
					log.Fatalln("Error marshalling the blocks: ", err)
				}
				w.Write(bs)
			}
			return
		} else if r.Method == "POST" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")

			var m map[string]string
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&m); err != nil {
				log.Fatalln("Error decoding the request body: ", err)
			}

			data := m["data"]
			blockchain.AddBlock(data)

			bs, err := json.Marshal(blockchain.Blocks)
			if err != nil {
				log.Fatalln("Error marshalling the blocks: ", err)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write(bs)
			return
		}
	})

	err := http.ListenAndServe(":9494", nil)
	if err != nil {
		log.Fatalln("Cannot start server: ", err)
	}
}
