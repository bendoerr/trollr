package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bendoerr/trollr/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostponeWrite", func() {
	It("should not write until the middleware returns", func() {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		c := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w_ http.ResponseWriter, r_ *http.Request) {
				next.ServeHTTP(w_, r_)
				Expect(w.Code).To(Equal(http.StatusOK)) // 200 is the default
				Expect(w.Body.Len()).To(BeZero())
			})
		}

		h := http.HandlerFunc(func(w_ http.ResponseWriter, r_ *http.Request) {
			w_.WriteHeader(http.StatusContinue)
			_, err := w_.Write([]byte{'f', 'o', 'o'})
			Expect(err).To(BeNil())
		})

		m := middleware.PostponeWriteMiddleware(c(h))
		m.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(http.StatusContinue))
		Expect(w.Body.Len()).To(Equal(3))
	})
})
