package web

import (
	"log"
	"net/http"
)

// middlewares
type reqResLogger struct {
	lgr *log.Logger
}

func newReqResLogger(lgr *log.Logger) http.Handler {
	return &reqResLogger{lgr}
}

func (lm *reqResLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lm.lgr.Printf("%s %s", r.Method, r.URL.Path)
	lm.lgr.Printf("%s", r.Header)
}
