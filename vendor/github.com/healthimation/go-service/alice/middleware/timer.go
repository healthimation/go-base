package middleware

import (
	"net/http"

	"github.com/justinas/alice"
)

// Timer can time a handler and log it
type Timer interface {
	Time(name string) alice.Constructor
}

type nullTimer struct{}

func (n *nullTimer) Time(name string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

func NewNullTimer() Timer {
	return &nullTimer{}
}
