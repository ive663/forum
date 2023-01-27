package delivery

import (
	"context"
	"net/http"

	"github.com/ive663/forum/internal/module"

	_ "github.com/mattn/go-sqlite3"
)

type key int

const (
	keyUserID key = iota
)

type (
	Middleware func(http.Handler) http.Handler
	Chain      []Middleware
)

func CreateChain(middlewares ...Middleware) Chain {
	var slice Chain
	return append(slice, middlewares...)
}

func (c Chain) Then(handler http.Handler) http.Handler {
	if handler == nil {
		handler = http.DefaultServeMux
	}
	for i := range c {
		handler = c[len(c)-1-i](handler)
	}
	return handler
}

func (h *Handler) authenticateUser(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := module.User{}
		c, err := r.Cookie("session")
		if err != nil {
			if err == http.ErrNoCookie {
				u.ID = 0
			}
			u.ID = 0
		} else {
			u.ID, err = h.services.GetUserIdByUUID(c.Value)
			if err != nil {
				u.ID = 0
			}
		}
		ctx := context.WithValue(r.Context(), keyUserID, u.ID)
		handler(w, r.WithContext(ctx))
	})
}
