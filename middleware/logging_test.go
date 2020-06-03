package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bendoerr/trollr/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var _ = Describe("LoggingMiddleware", func() {
	var (
		m            http.Handler
		r            *http.Request
		w            *httptest.ResponseRecorder
		observedLogs *observer.ObservedLogs
	)

	BeforeEach(func() {
		var core zapcore.Core
		core, observedLogs = observer.New(zapcore.InfoLevel)
		logger := zap.New(core)

		m = middleware.LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte{'h', 'e', 'l', 'l', 'o'})
			Expect(err).To(BeNil())
		}), logger)
		r = httptest.NewRequest("GET", "/", nil)
		w = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		m.ServeHTTP(w, r)
	})

	It("should log standard request attributes", func() {
		logs := observedLogs.FilterMessage(r.RequestURI).TakeAll()
		Expect(logs).To(HaveLen(1))
		Expect(logs[0].Message).To(Equal(r.RequestURI))
		Expect(logs[0].Context).To(ContainElement(zap.Int64("request.content_length", r.ContentLength)))
		Expect(logs[0].Context).To(ContainElement(zap.String("request.host", r.Host)))
		Expect(logs[0].Context).To(ContainElement(zap.String("request.method", r.Method)))
		Expect(logs[0].Context).To(ContainElement(zap.String("request.user_agent", r.UserAgent())))
	})

	When("HTTP request headers are sent", func() {
		BeforeEach(func() {
			r.Header.Set("Test-Header", "Test-Value")
		})

		It("should log HTTP request headers", func() {
			logs := observedLogs.FilterMessage(r.RequestURI).TakeAll()
			Expect(logs).To(HaveLen(1))
			Expect(logs[0].Context).To(ContainElement(zap.String("request.header.test_header", "Test-Value")))
		})
	})

	It("should log HTTP status code & content length", func() {
		logs := observedLogs.FilterMessage(r.RequestURI).TakeAll()
		Expect(logs).To(HaveLen(1))
		Expect(logs[0].Context).To(ContainElement(zap.Int("response.code", 200)))
		Expect(logs[0].Context).To(ContainElement(zap.Int("response.content-length", 5)))
	})

	When("the same HTTP request header is repeated", func() {
		BeforeEach(func() {
			r.Header.Add("Test-Header", "Test-Value-1")
			r.Header.Add("Test-Header", "Test-Value-2")
			r.Header.Add("Test-Header", "Test-Value-3")
		})

		It("should log HTTP each repeated request header", func() {
			logs := observedLogs.FilterMessage(r.RequestURI).TakeAll()
			Expect(logs).To(HaveLen(1))
			Expect(logs[0].Context).To(ContainElement(zap.String("request.header.test_header_0", "Test-Value-1")))
			Expect(logs[0].Context).To(ContainElement(zap.String("request.header.test_header_1", "Test-Value-2")))
			Expect(logs[0].Context).To(ContainElement(zap.String("request.header.test_header_2", "Test-Value-3")))
		})
	})
})
