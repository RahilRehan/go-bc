package web

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	gobc "github.com/RahilRehan/go-bc"
)

type walletHandler struct {
	wallets []gobc.Wallet
}

func (t *walletHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t.handleGet(w, r)
	case "POST":
		t.handlePost(w, r)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (t *walletHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")
}

func (t *walletHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST")
}

func NewWebServer(port string) {

	wh := newWalletHandler()
	mh := newMiddlewareHandler(wh)
	mh.Use(newReqResLogger(log.New(os.Stdout, "", log.LstdFlags)))

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	mux.Handle("/wallet", mh)

	server := &http.Server{
		Addr:           port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("=====================================")
	log.Println("Web Server started on port " + port)
	fmt.Println("=====================================")
	log.Fatalln(server.ListenAndServe())

}

func newWalletHandler() http.Handler {
	return &walletHandler{
		wallets: make([]gobc.Wallet, 0),
	}
}
