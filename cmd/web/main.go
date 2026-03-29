package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Emin-07/Sentinel/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
}
type application struct {
	errorlog  *log.Logger
	infoLog   *log.Logger
	processes *models.ProcessModel
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	dsn := flag.String("dsn", "web:pass@/sentinel?parseTime=true", "MySQL data source name")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.LUTC|log.Ldate|log.Ltime)
	port, err := strconv.Atoi(cfg.addr[1:])
	if port <= 1024 {
		errorLog.Print("Bind: permission denied")
		return
	}
	if err != nil {
		errorLog.Print(err)
	}
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	app := &application{
		errorlog:  errorLog,
		infoLog:   infoLog,
		processes: &models.ProcessModel{DB: db},
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(&cfg),
	}
	app.infoLog.Printf("Starting a server on %s", cfg.addr)

	err = srv.ListenAndServe()
	if err != nil {
		app.errorlog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// if process is in active for a day, then liquadate it
