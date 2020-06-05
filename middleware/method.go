package middleware

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func PostMethodOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = jsoniter.NewEncoder(w).Encode(struct {
				Error string
			}{
				Error: http.StatusText(http.StatusMethodNotAllowed) + ", expect POST",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
