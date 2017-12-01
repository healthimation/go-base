package test

import (
	"net/http"

	"github.com/healthimation/go-service/alice/middleware"
	"github.com/healthimation/go-service/service"
)

type pong struct {
	Message string `json:"message"`
}

// Ping returns a ping handler
func Ping() http.Handler {
	return middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ret := pong{Message: "pong!"}
		return service.WriteJSONResponse(w, http.StatusOK, ret)
	})
}
