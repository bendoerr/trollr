package middleware

import (
	"net/http"

	servertiming "github.com/mitchellh/go-server-timing"
)

func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timing := servertiming.FromContext(r.Context())
		metric := timing.NewMetric("request").Start()
		next.ServeHTTP(w, r)
		metric.Stop()

		// The servertiming.Middleware probably has already tried to write out this header.
		// We can do it here as long as we trust that PostponeWriteMiddleware is in our chain.
		w.Header().Set(servertiming.HeaderKey, timing.String())
	})
}
