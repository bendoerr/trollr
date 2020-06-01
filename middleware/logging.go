package middleware

import (
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	servertiming "github.com/mitchellh/go-server-timing"
	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs := []zap.Field{
			zap.Int64("request.request_uri", r.ContentLength),
			zap.String("request.host", r.Host),
			zap.String("request.method", r.Method),
			zap.String("request.proto", r.Proto),
			zap.String("request.referer", r.Referer()),
			zap.String("request.remote_addr", r.RemoteAddr),
			zap.String("request.request_uri", r.RequestURI),
			zap.String("request.user_agent", r.UserAgent()),
		}

		for k, v := range r.Header {
			if len(v) > 1 {
				for i := range v {
					fs = append(fs, zap.String("request.header."+strings.ReplaceAll(strings.ToLower(k), "-", "_")+"_"+string(i), v[i]))
				}
			} else {
				fs = append(fs, zap.String("request.header."+strings.ReplaceAll(strings.ToLower(k), "-", "_"), v[0]))
			}
		}

		var returnCode int
		var contentLength int

		next.ServeHTTP(httpsnoop.Wrap(w, httpsnoop.Hooks{
			WriteHeader: func(original httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				return func(code int) {
					returnCode = code
					original(code)
				}
			},
			Write: func(original httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				return func(b []byte) (int, error) {
					i, err := original(b)
					contentLength = contentLength + i
					return i, err
				}
			},
		}), r)

		fs = append(fs, zap.Int("response.code", returnCode))
		fs = append(fs, zap.Int("response.content-length", contentLength))

		timing := servertiming.FromContext(r.Context())
		if timing != nil {
			for i := range timing.Metrics {
				m := timing.Metrics[i]
				fs = append(fs, zap.Duration("timing."+strings.ReplaceAll(strings.ToLower(m.Name), "-", "_"), m.Duration))
			}
		}

		logger.Info(r.URL.String(), fs...)
	})
}
