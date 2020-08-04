package middleware

import "net/http"

// NoCache will add the following headers to the response:
// 	Cache-Control: no-cache, no-store, must-revalidate
// 	Pragma: no-cache
// 	Expires: 0
func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate") // HTTP 1.1
		w.Header().Set("Pragma", "no-cache")                                              // HTTP 1.0
		w.Header().Set("Expires", "0")                                                    // Proxies
		h.ServeHTTP(w, r)
	})
}
