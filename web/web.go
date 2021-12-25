package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	gobc "github.com/RahilRehan/go-bc"
)

type gobcHandler struct {
	wallets []gobc.Wallet
}

type txRequest struct {
	Amount int64 `json:"amount"`
}

func (gh *gobcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		gh.handleGet(w, r)
	case "POST":
		gh.handlePost(w, r)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (gh *gobcHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://localhost:8080/transaction")
	if err != nil {
		log.Fatalln("can't get transactions from server", err)
	}
	var txPool gobc.TransactionPool
	bs, _ := io.ReadAll(res.Body)
	json.Unmarshal(bs, &txPool)
	if err != nil {
		log.Fatalln("can't decode transactions from server", err)
	}
	// txPool := gobc.TransactionPool{Transactions: txs}
	fmt.Fprintf(w, "Transactions: %v", txPool)
}

func (gh *gobcHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	sender := gobc.NewWallet()
	receiver := gobc.NewWallet()
	gh.wallets = append(gh.wallets, sender)
	gh.wallets = append(gh.wallets, receiver)

	if strings.HasPrefix(r.URL.Path, "/gobc/transactions") {
		var txReq txRequest
		bs, _ := io.ReadAll(r.Body)
		json.Unmarshal(bs, &txReq)
		tx := gobc.NewTransaction(&sender, &receiver, txReq.Amount)
		bs, err := json.Marshal(tx)
		if err != nil {
			log.Fatalln("cannot marshal transaction")
		}

		_, err = http.Post("http://localhost:8080/transaction", "application/json", bytes.NewBuffer(bs))
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/gobc/transactions", http.StatusSeeOther)
	}
}

func NewWebServer(port string) {

	gh := newGobcHandler()
	mh := newMiddlewareHandler(gh)
	mh.Use(newReqResLogger(log.New(os.Stdout, "", log.LstdFlags)))

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	mux.Handle("/gobc/", mh)
	mux.Handle("/", http.FileServer(http.Dir("./web/static")))

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

func newGobcHandler() http.Handler {
	return &gobcHandler{
		wallets: make([]gobc.Wallet, 0),
	}
}
