package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func PingWebsite(url string) error {
	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Head(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (app *application) cleanProcesses(w http.ResponseWriter) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		processes, err := app.processes.Latest()
		if err != nil {
			app.errorlog.Print(err)
		}
		json.NewEncoder(w).Encode(map[string]bool{"WindowChanged": false})
		for _, process := range processes {
			if err := PingWebsite(process.Title); err != nil {
				process.Active = false
			}
			if !process.Active && time.Since(process.StartedAt).Minutes() >= 1 { //change minutes to hours, and check every 5 mins or so
				app.processes.Delete(process.ID)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]bool{"WindowChanged": true})
			}
		}
	}
}

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
	go app.cleanProcesses(w)
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
	link := r.Response.Body
	fmt.Fprintf(w, "%v", link)
}
