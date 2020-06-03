package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bendoerr/trollr/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recovery", func() {
	It("should recover a panic", func() {
		m := middleware.Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("")
		}))
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		Expect(func() { m.ServeHTTP(w, r) }).NotTo(Panic())
		Expect(w.Result().StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
