package main

import (
	"net/http"
	"text/template"
	"time"
)

type Process struct {
	ID        int
	Title     string
	Observing time.Time
	Active    bool
}
type testStruct struct {
	Processes []Process
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	process := &Process{
		ID:        1,
		Title:     "title",
		Observing: time.Now(),
		Active:    true,
	}
	test := &testStruct{
		Processes: []Process{*process},
	}
	if r.URL.Path != "/" {
		app.notFound(w)
	}

	if r.Method != "GET" {
		app.clientError(w, http.StatusMethodNotAllowed)
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", test)
	if err != nil {
		app.serverError(w, err)
	}
}

// func (app *application) viewProcess(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if id >= 1 || err != nil {
// 		app.serverError(w, err)
// 		return
// 	}

// 	if r.Method != "GET" {
// 		app.clientError(w, http.StatusMethodNotAllowed)
// 	}

// 	files := []string{
// 		"./ui/html/base.html",
// 		"./ui/html/pages/home.html",
// 		"./ui/html/partials/nav.html",
// 	}

// 	ts, err := template.ParseFiles(files...)
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}
// 	err = ts.ExecuteTemplate(w, "base", nil)
// 	if err != nil {
// 		app.serverError(w, err)
// 	}
// }
