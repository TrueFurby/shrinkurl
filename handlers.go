package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(app *App) *mux.Router {
	var r = mux.NewRouter()

	r.Path("/{hash:[a-z0-9]{6}}").Methods("GET").
		HandlerFunc(app.redirectHandler)

	r.Path("/add").Methods("POST").
		HandlerFunc(app.addHandler)

	r.Path("/check").Methods("GET").
		HandlerFunc(app.checkHandler)

	r.Path("/remove").Methods("DELETE").
		HandlerFunc(app.removeHandler)

	return r
}

func (a *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	var hash = mux.Vars(r)["hash"]

	if cached := a.Cache.Get(hash); cached != "" {
		log.Println("using cached url", cached)
		http.Redirect(w, r, cached, http.StatusMovedPermanently)
		return
	}

	u, err := a.UrlStore().Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u.Id == 0 {
		http.NotFound(w, r)
		return
	}

	a.Cache.Add(hash, u.Url)

	log.Println(hash, "redirecting to", u.Url)
	http.Redirect(w, r, u.Url, http.StatusMovedPermanently)
}

func (a *App) addHandler(w http.ResponseWriter, r *http.Request) {
	var inputUrl = r.FormValue("url")

	destUrl, err := parseUrl(inputUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := a.UrlStore().GetByUrl(destUrl.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u.Id == 0 {
		u.Hash = makeHash(destUrl)
		u.Url = destUrl.String()
		if err := a.UrlStore().Update(&u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(destUrl, "added")
		writeJson(w, http.StatusCreated, u)
		return
	}

	writeJson(w, http.StatusOK, u)
}

func (a *App) checkHandler(w http.ResponseWriter, r *http.Request) {
	var hash = r.FormValue("hash")

	u, err := a.UrlStore().Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u.Id == 0 {
		http.NotFound(w, r)
		return
	}

	writeJson(w, http.StatusOK, u)
}

func (a *App) removeHandler(w http.ResponseWriter, r *http.Request) {
	var hash = r.FormValue("hash")

	if u, err := a.UrlStore().Get(hash); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u.Id == 0 {
		http.NotFound(w, r)
		return
	}

	err := a.UrlStore().Remove(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.Cache.Remove(hash)

	log.Println(hash, "removed")
	writeJson(w, http.StatusNoContent, nil)
}

func writeJson(w http.ResponseWriter, status int, a interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if a != nil {
		if err := json.NewEncoder(w).Encode(a); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
