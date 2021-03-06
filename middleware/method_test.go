package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bendoerr/trollr/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostMethodOnlyMiddleware", func() {

	var (
		m http.Handler
		r *http.Request
		w *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		m = middleware.PostMethodOnlyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte{})
			Expect(err).To(BeNil())
		}))
		w = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		m.ServeHTTP(w, r)
	})

	When("the HTTP method is POST", func() {
		BeforeEach(func() {
			r = httptest.NewRequest("POST", "/", nil)
		})
		It("should NOT set StatusMethodNotAllowed", func() {
			Expect(w.Result().StatusCode).To(Not(Equal(http.StatusMethodNotAllowed)))
		})
	})

	When("the HTTP method is PUT", func() {
		BeforeEach(func() {
			r = httptest.NewRequest("PUT", "/", nil)
		})
		It("should set StatusMethodNotAllowed", func() {
			Expect(w.Result().StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	When("the HTTP method is GET", func() {
		BeforeEach(func() {
			r = httptest.NewRequest("GET", "/", nil)
		})
		It("should set StatusMethodNotAllowed", func() {
			Expect(w.Result().StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})
	})
})
