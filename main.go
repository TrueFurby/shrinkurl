package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Storage
	*mux.Router
	Cache *Cache
}

func NewApp() *App {
	var a = &App{
		Storage: NewStorage("sqlite3", "./url.db"),
		Cache:   NewCache(10),
	}
	a.Router = NewRouter(a)
	return a
}

func main() {
	var app = NewApp()

	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	}).Methods("GET")

	if err := http.ListenAndServe(":7788", app); err != nil {
		log.Fatal(err)
	}
}
