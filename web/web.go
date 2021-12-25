package web

import (
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
	txPool  gobc.TransactionPool
}

func (t *gobcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (t *gobcHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/gobc/transactions") {
		fmt.Fprintln(w, "YO! transactions!")
		for _, tx := range t.txPool.Transactions {
			fmt.Fprintln(w, tx)
		}
		return
	}
}

func (t *gobcHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	sender := gobc.NewWallet()
	receiver := gobc.NewWallet()
	t.wallets = append(t.wallets, sender)
	t.wallets = append(t.wallets, receiver)

	var amount int64
	if strings.HasPrefix(r.URL.Path, "/gobc/transactions") {
		bs, _ := io.ReadAll(r.Body)
		json.Unmarshal(bs, &amount)
		tx := gobc.NewTransaction(&sender, &receiver, amount)
		t.txPool.Add(tx)
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
