package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to GO Blockchain!"))
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Cannot start server: ", err)
	}
}
