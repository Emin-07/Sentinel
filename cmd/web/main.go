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

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.LUTC|log.Ldate|log.Ltime)

	app := &application{
		errorlog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(&cfg),
	}
	app.infoLog.Printf("Starting a server on %s", cfg.addr)

	err := srv.ListenAndServe()
	if err != nil {
		app.errorlog.Fatal(err)
	}
}
