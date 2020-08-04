package chain

import (
	"fmt"
	"log"
	"net/http"

	"github.com/healthimation/go-service/service"
)

// NotFoundHandler is the generic not found handler to return a problem
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not Found: %s", r.URL.String())
	service.WriteProblem(w, "No route found", "ERROR_NOT_FOUND", http.StatusNotFound)
}

// MethodNotAllowedHandler returns a method not allowed handler - Implements vestigo.MethodNotAllowedHandlerFunc
func MethodNotAllowedHandler(method string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method (%s) not allowed on url (%s)", r.Method, r.URL.String())
		service.WriteProblem(w, fmt.Sprintf("Method not allowed.  Allowed methods are (%s)", method), "ERROR_METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
	}
}
