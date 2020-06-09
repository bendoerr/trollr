// Trollr: A HTTP/JSON API for Troll
//
// Trollr is a simple wrapper around the amazing Troll: A dice roll language and calculator created by Torben Mogensen.
// The wrapper simply exposes and HTTP/JSON server that executes Troll, parses the results and returns it. The server
// has some built in rate-limiting and pooling to prevent abuse. I created this small server to support a Discord bot
// that I am working on.
//
//     Schemes: https
//     Host: trollr.live
//     BasePath: /api
//     Version: 0.1.0-alpha
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Ben Doerr <craftsman@bendoer.me> https://trollr.live
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/markusthoemmes/goautoneg"

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
	swaggerFile string
	swaggerRedirect string
}

func NewAPI(listen string, troll *Troll, logger *zap.Logger, swaggerFile, swaggerRedirect string) *API {
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
	api.mux.Handle("/", http.HandlerFunc(api.DefaultNotFound))
	api.mux.Handle("/swagger.json", http.HandlerFunc(api.SwaggerJson))
	api.mux.Handle("/swagger", http.HandlerFunc(api.SwaggerRedirect))
	api.mux.Handle("/roll", middleware.PostMethodOnlyMiddleware(http.HandlerFunc(api.Roll)))
	api.mux.Handle("/calc", middleware.PostMethodOnlyMiddleware(http.HandlerFunc(api.Calc)))

	h := http.Handler(api.mux)
	h = middleware.JsonContentTypeMiddleware(h)
	h = middleware.NoticeHeadersMiddleware(h)
	h = middleware.TimingMiddleware(h)
	h = middleware.LoggingMiddleware(h, logger)
	h = servertiming.Middleware(h, nil)
	h = middleware.PostponeWriteMiddleware(h)
	h = middleware.Recovery(h)
	h = tollbooth.LimitHandler(tollbooth.NewLimiter(1, nil), h)

	api.server.Handler = h

	api.swaggerFile = swaggerFile
	api.swaggerRedirect = swaggerRedirect

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

func (api *API) DefaultNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = jsoniter.NewEncoder(w).Encode(RollsResult{
		Error: http.StatusText(http.StatusNotFound),
	})
}

func (api *API) SwaggerJson(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, api.swaggerFile)
}

func (api *API) SwaggerRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, api.swaggerRedirect, http.StatusTemporaryRedirect)
}

// swagger:operation POST /roll API roll
//
// Roll Dice
//
// Given a roll definition this will delegate the roll to Troll and return the
// results structured as JSON.
//
// ---
// consumes:
// - text/plain
// produces:
// - application/json
// parameters:
// - name: "d"
//   in: query
//   description: "The Troll roll definition. This can passed as the query parameter 'd' or in the request body."
//   type: "string"
//   required: false
// - name: "d"
//   in: body
//   description: "The Troll roll definition. This can passed as the query parameter 'd' or in the request body."
//   schema:
//     type: string
// - name: "n"
//   in: "query"
//   description: "The number of times to repeat the roll"
//   type: "integer"
//   required: false
// responses:
//   '200':
//     description: "The results from rolling the dice"
//     schema:
//       "$ref": "#/definitions/RollsResult"
//   '400':
//     description: "The error will be populated in the result"
//     schema:
//       "$ref": "#/definitions/RollsResult"
func (api *API) Roll(w http.ResponseWriter, r *http.Request) {
	d := r.URL.Query().Get("d")
	n := r.URL.Query().Get("n")

	if len(d) < 1 {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
		}
		body := strings.TrimSpace(string(bytes))
		if len(body) > 0 {
			d = body
		}
	}

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

	accepts := goautoneg.ParseAccept(r.Header.Get("Accept"))
	if len(accepts) > 0 && accepts[0].Type == "text" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte(res.RollsRaw))
		fmt.Println(err)
	} else {
		err := jsoniter.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// swagger:operation POST /calc API calc
//
// Calculate the probabilities of dice roll.
//
// Given a roll definition this will delegate the roll to Troll and return the
// probabilities structured as JSON.
//
// ---
// consumes:
// - text/plain
// produces:
// - application/json
// parameters:
// - name: "d"
//   in: query
//   description: "The Troll roll definition. This can passed as the query parameter 'd' or in the request body."
//   type: "string"
//   required: false
// - name: "d"
//   in: body
//   description: "The Troll roll definition. This can passed as the query parameter 'd' or in the request body."
//   schema:
//     type: string
// - name: "c"
//   in: "query"
//   description: "What kind of cumulative probabilities you would like. One of 'ge' (default), 'gt', 'le', or 'lt'."
//   type: "string"
//   required: false
// responses:
//   '200':
//     description: "The probabilities of rolling the dice"
//     schema:
//       "$ref": "#/definitions/CalcResult"
//   '400':
//     description: "The error will be populated in the result"
//     schema:
//       "$ref": "#/definitions/CalcResult"
func (api *API) Calc(w http.ResponseWriter, r *http.Request) {
	d := r.URL.Query().Get("d")
	c := r.URL.Query().Get("c")

	if len(d) < 1 {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
		}
		body := strings.TrimSpace(string(bytes))
		if len(body) > 0 {
			d = body
		}
	}

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

	accepts := goautoneg.ParseAccept(r.Header.Get("Accept"))
	if len(accepts) > 0 && accepts[0].Type == "text" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte(res.ProbabilitiesRaw))
		fmt.Println(err)
	} else {
		err := jsoniter.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}
