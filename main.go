package main

import (
	"log"
	"net/http"
)

type App struct {
	*Storage
	*http.ServeMux
	cache *Cache
}

func main() {
	var app = &App{
		Storage:  NewStorage("./url.db"),
		ServeMux: http.NewServeMux(),
		cache:    NewCache(10),
	}

	app.HandleFunc("/add", app.addHandler)
	app.HandleFunc("/check", app.checkHandler)
	app.HandleFunc("/remove", app.removeHandler)
	app.HandleFunc("/", app.indexHandler)

	if err := http.ListenAndServe(":7788", app); err != nil {
		log.Fatal(err)
	}
}
