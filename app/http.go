package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bendoerr/trollr/middleware"
	"github.com/didip/tollbooth"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	servertiming "github.com/mitchellh/go-server-timing"
	"go.uber.org/zap"
)

type API struct {
	server *http.Server
	mux    *http.ServeMux
	troll  *Troll
}

func NewAPI(listen string, troll *Troll, logger *zap.Logger) *API {
	api := &API{
		troll: troll,
	}

	api.server = &http.Server{
		Addr:              listen,
		Handler:           nil,
		TLSConfig:         nil,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	// Have jsoniter handle naming
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
	extra.RegisterFuzzyDecoders()

	api.mux = &http.ServeMux{}

	api.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = jsoniter.NewEncoder(w).Encode(RollsResult{
			Error: http.StatusText(http.StatusNotFound),
		})
	}))

	api.mux.Handle("/roll", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := r.URL.Query().Get("d")
		n := r.URL.Query().Get("n")

		if len(d) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_ = jsoniter.NewEncoder(w).Encode(RollsResult{
				Error: "no roll definition given",
			})
			return
		}

		c := 1
		if len(n) > 0 {
			i, err := strconv.Atoi(n)
			if err == nil {
				c = i
			}
		}

		res := api.troll.MakeRolls(r.Context(), c, d)
		if res.Err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err := jsoniter.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}))

	api.mux.Handle("/calc", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := r.URL.Query().Get("d")
		c := r.URL.Query().Get("c")

		if len(d) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_ = jsoniter.NewEncoder(w).Encode(RollsResult{
				Error: "no roll definition given",
			})
			return
		}

		res := api.troll.CalcRoll(r.Context(), d, c)
		if res.Err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err := jsoniter.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}))

	h := http.Handler(api.mux)
	h = middleware.GetMethodOnlyMiddleware(h)
	h = middleware.JsonContentTypeMiddleware(h)
	h = middleware.NoticeHeadersMiddleware(h)
	h = middleware.TimingMiddleware(h)
	h = middleware.LoggingMiddleware(h, logger)
	h = servertiming.Middleware(h, nil)
	h = middleware.PostponeWriteMiddleware(h)
	h = middleware.Recovery(h)
	h = tollbooth.LimitHandler(tollbooth.NewLimiter(1, nil), h)
	api.server.Handler = h

	return api
}

func (api *API) Start() {
	go func() {
		err := api.server.ListenAndServe()
		fmt.Println("server exited: ", err)
	}()
}

func (api *API) Stop() error {
	return api.server.Shutdown(context.Background())
}
