package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/bendoerr/trollr/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NoticeHeadersMiddleware", func() {
	It("should set the Notice headers", func() {
		m := middleware.NoticeHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)

		Expect(w.Result().Header).To(HaveKey("Notice-Message"))
		Expect(w.Result().Header).To(HaveKey("Notice-Troll-Author"))
		Expect(w.Result().Header).To(HaveKey("Notice-Troll-Url"))
		Expect(w.Result().Header).To(HaveKey("Notice-Troll-Manual"))
	})
})
