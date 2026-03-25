package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}
type application struct {
	errorlog *log.Logger
	infoLog  *log.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":http", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.LUTC|log.Ldate|log.Ltime)

	app := &application{
		errorlog: errorLog,
		infoLog:  infoLog,
	}
	app.infoLog.Printf("Starting a server on %s", cfg.addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	err := http.ListenAndServe(cfg.addr, mux)
	if err != nil {
		app.errorlog.Fatal(err)
	}
}
