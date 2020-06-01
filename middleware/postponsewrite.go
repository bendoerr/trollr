package middleware

import (
	"fmt"
	"net/http"

	"github.com/felixge/httpsnoop"
)

func PostponeWriteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var delayedHeaders []func()
		var delayedWrites []func()

		next.ServeHTTP(httpsnoop.Wrap(w, httpsnoop.Hooks{
			Write: func(original httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				return func(b []byte) (int, error) {
					delayedWrites = append(delayedWrites, func() {
						_, err := original(b)
						if err != nil {
							fmt.Println("write failure: ", err)
						}
					})
					return len(b), nil
				}
			},
			WriteHeader: func(original httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				return func(code int) {
					delayedHeaders = append(delayedHeaders, func() {
						original(code)
					})
				}
			},
		}), r)

		// Finally write the body (thus closing the response for any other modifications
		for _, hf := range delayedHeaders {
			hf()
		}
		for _, wf := range delayedWrites {
			wf()
		}
	})
}
