package middleware

import (
	"context"
	"net/http"

	"github.com/pborman/uuid"
)

const (
	headerRequestID = "x-request-id"
)

// RequestID adds a request id to the context, if available it will use the one in the x-request-id header
func RequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check header for existing request id
		rID := r.Header.Get(headerRequestID)
		if len(rID) == 0 {
			//generate id
			rID = uuid.New()
		}

		// add rID to the context
		ctx := context.WithValue(r.Context(), contextKeyRequestID, rID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
