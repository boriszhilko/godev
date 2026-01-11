package middleware

import (
	"net/http"
	"strings"
)

const apiKeyHeader = "X-API-Key"

// Auth validates API keys from the request header.
// validKeys is a list of accepted API keys.
func Auth(validKeys []string) func(http.Handler) http.Handler {
	keyMap := make(map[string]bool)
	for _, key := range validKeys {
		keyMap[key] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := strings.TrimSpace(r.Header.Get(apiKeyHeader))

			if apiKey == "" || !keyMap[apiKey] {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"success":false,"error":"Invalid or missing API key","code":"UNAUTHORIZED"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
