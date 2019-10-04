package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"net/url"
)

type App struct {
	Router      *mux.Router
	MiddleWares *MiddleWare
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortLinkResp struct {
	ShortLink string `json:"short_link"`
}

// initialize

func (a *App) Initialize() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.MiddleWares = &MiddleWare{}
	a.initializeRoutes()

}
func (a *App) initializeRoutes() {
	var (
		m alice.Chain
	)
	m = alice.New(a.MiddleWares.LoggingHandler, a.MiddleWares.RecoverHandler)
	// 加入 middleware
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortLink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortLink)).Methods("GET")
	a.Router.Handle("/{shortLink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")

}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var (
		req shortenReq
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		respondWithError(w, NewStatusError(fmt.Errorf("parse params failed %v", r.Body)))
		return
	}
	if err = validator.Validate(req); err != nil {
		fmt.Println(err)
		respondWithError(w, NewStatusError(fmt.Errorf("validate params failed %v", req)))
		return
	}

	defer r.Body.Close()

	fmt.Println(req)

}

func (a *App) getShortLink(w http.ResponseWriter, r *http.Request) {
	var (
		vals url.Values
		s    string
	)

	vals = r.URL.Query()
	s = vals.Get("shortLink")

	fmt.Println(s)

	panic(s)

}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	var (
		vars map[string]string
	)
	vars = mux.Vars(r)
	fmt.Println(vars["shortLink"])

}

// run and start listening
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))

}

func respondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {

	case Error:
		log.Println("HTTP", e.Status(), e.Error())
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	var (
		resp []byte
		err  error
	)
	if resp, err = json.Marshal(payload); err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(resp)

}
