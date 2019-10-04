package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"net/url"
)

type App struct {
	Router *mux.Router
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
	a.initializeRoutes()

}
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortLink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortLink).Methods("GET")
	a.Router.HandleFunc("/{shortLink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")

}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var (
		req shortenReq
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err)
		respondWithError(w,StatusError{
			http.StatusBadRequest,
			fmt.Errorf("parse params failed %s",r.Body)
		})
		return
	}
	if err = validator.Validate(req); err != nil {
		fmt.Println(err)
		respondWithError(w,StatusError{
			http.StatusBadRequest,
			fmt.Errorf("validate params failed %s",r.Body)
		})
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

func respondWithError(writer http.ResponseWriter, err error ) {
	switch e := err.(type) {

	case Error:
		log.Println("HTTP",e.Status(),e)
		respondWithJSON()


	
	}
}
