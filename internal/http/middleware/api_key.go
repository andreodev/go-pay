package middleware

import (
	"net/http"
)

func APIKeyAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestAPIKey := r.Header.Get("X-API-Key")

			if requestAPIKey == "" {
				http.Error(w, "missing api key", http.StatusUnauthorized)
				return
			}

			if requestAPIKey != apiKey {
				http.Error(w, "invalid api key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
