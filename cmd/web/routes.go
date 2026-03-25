package main

import "net/http"

func (app *application) routes(cfg *config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/process/view", app.processView)
	return mux
}
