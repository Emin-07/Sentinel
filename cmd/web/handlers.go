package main

import (
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
	}

	if r.Method != "GET" {
		app.clientError(w, http.StatusMethodNotAllowed)
	}
	processes, err := app.processes.Latest()
	if err != nil {
		app.errorlog.Print(err)
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
	data := &templateData{Processes: processes}
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) viewProcess(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if id < 1 {
		app.notFound(w)
		return
	}
	if r.Method != "GET" {
		app.clientError(w, http.StatusMethodNotAllowed)
	}
	process, err := app.processes.Get(id)
	if err != nil {
		app.errorlog.Print(err)
	}
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/view.html",
		"./ui/html/partials/nav.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Process: process}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	process, err := app.processes.Insert()
}
