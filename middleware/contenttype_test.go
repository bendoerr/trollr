package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bendoerr/trollr/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JsonContentTypeMiddleware", func() {
	It("should set the Content-Type header", func() {
		m := middleware.JsonContentTypeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)

		Expect(w.Header()).To(HaveKeyWithValue("Content-Type", []string{"application/json; charset=utf-8"}))
	})
})
