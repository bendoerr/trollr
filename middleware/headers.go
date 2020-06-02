package middleware

import "net/http"

func NoticeHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Notice-Message", "The 'Trollr' API is a simple server that wraps the amazing 'Troll' program. This API is not associated with the author of the 'Troll' program.")
		w.Header().Set("Notice-Troll-Author", "Torben Mogensen <torbenm@di.ku.dk>")
		w.Header().Set("Notice-Troll-Url", "http://hjemmesider.diku.dk/~torbenm/Troll/")
		w.Header().Set("Notice-Troll-Manual;", "http://hjemmesider.diku.dk/~torbenm/Troll/manual.pdf")
		next.ServeHTTP(w, r)
	})
}
