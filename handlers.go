package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path

	if path == "/" {
		http.ServeFile(w, r, "./index.html")
	} else {
		a.redirectHandler(w, r)
	}
}

func (a *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	var hash = r.URL.Path[1:]

	if cached := a.cache.Get(hash); cached != "" {
		log.Println("using cached url", cached)
		http.Redirect(w, r, cached, http.StatusMovedPermanently)
		return
	}

	u, err := a.Url.Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u.Id == 0 {
		http.NotFound(w, r)
		return
	}

	a.cache.Add(hash, u.Url)

	log.Println(hash, "redirecting to", u.Url)
	http.Redirect(w, r, u.Url, http.StatusMovedPermanently)
}

func (a *App) addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}
	var inputUrl = r.FormValue("url")

	destUrl, err := parseUrl(inputUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := a.Url.GetByUrl(destUrl.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u.Id == 0 {
		u.Hash = makeHash(destUrl)
		u.Url = destUrl.String()
		if err := a.Url.Update(&u); err != nil {
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
	if r.Method != "GET" {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}
	var hash = r.FormValue("hash")

	u, err := a.Url.Get(hash)
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
	if r.Method != "DELETE" {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}
	var hash = r.FormValue("hash")

	if u, err := a.Url.Get(hash); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u.Id == 0 {
		http.NotFound(w, r)
		return
	}

	err := a.Url.Remove(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.cache.Remove(hash)

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
