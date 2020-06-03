package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bendoerr/trollr/middleware"
	servertiming "github.com/mitchellh/go-server-timing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Timing", func() {
	It("should set a server timing header", func() {
		m := middleware.TimingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
		}))
		m = servertiming.Middleware(m, nil)
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		Expect(w.Result().Header).To(HaveKey(servertiming.HeaderKey))
	})
	It("should handle a nil servertiming", func() {
		m := middleware.TimingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
		}))
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		Expect(func() { m.ServeHTTP(w, r) }).ToNot(Panic())
		Expect(w.Result().Header).ToNot(HaveKey(servertiming.HeaderKey))
	})
})
