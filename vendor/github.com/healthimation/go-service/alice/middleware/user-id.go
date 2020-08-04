package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserIDInjector add the user id to the context this is not validated for perf reasons
func UserIDInjector(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		t := r.Header.Get("Authorization")
		t = strings.TrimPrefix(t, "Bearer ")
		parts := strings.Split(t, ".")

		if len(parts) == 3 {
			by, err := jwt.DecodeSegment(parts[1])
			if err == nil {
				claims := make(map[string]interface{})
				err = json.Unmarshal(by, &claims)
				if err == nil {
					if userIDi, ok := claims["sub"]; ok {
						if v, ok := userIDi.(string); ok {
							userID = v
						}
					}

				}
			}
		}

		// add userID to the context
		ctx := context.WithValue(r.Context(), contextKeyInsecureUserID, userID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
