package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.amagkn.ru/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	addr, dsn := parseCommandArgs()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("starting server", "addr", addr)

	err = http.ListenAndServe(addr, app.routes())

	if err != nil {
		logger.Error(err.Error())
	}
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func parseCommandArgs() (string, string) {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:web@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	return *addr, *dsn
}
