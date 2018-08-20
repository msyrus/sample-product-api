package middleware

import (
	"net/http"
)

// Auth middleware checks API authorization
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		// NOTE: token validation goes here
		// currently any non empty token is valid
		if tok == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
