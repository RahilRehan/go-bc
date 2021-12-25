package web

import (
	"net/http"
)

// middlewareHandler is a http.Handler that applies middlewares to the request
// approach1: list of middlewares and main middleware (implemented)
// approach2: each middleware handler has a next pointer which points to the next middleware handler => we can do pre and post middleware logic
// approach 2 is better because we can process the request and pass it to next middlewares, extract the resp from prev middleware and process it.
type middlewareHandler struct {
	middlewares []http.Handler
	next        http.Handler
}

func (m *middlewareHandler) Use(f http.Handler) {
	m.middlewares = append(m.middlewares, f)
}

func (m *middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, f := range m.middlewares {
		f.ServeHTTP(w, r)
	}
	m.next.ServeHTTP(w, r)
}

func newMiddlewareHandler(next http.Handler) *middlewareHandler {
	return &middlewareHandler{
		middlewares: []http.Handler{},
		next:        next,
	}
}
