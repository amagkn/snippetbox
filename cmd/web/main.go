package main

import (
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := parseCommandArgs()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := application{
		logger: logger,
	}

	logger.Info("starting server", "addr", addr)

	err := http.ListenAndServe(addr, app.routes())

	if err != nil {
		logger.Error(err.Error())
	}
	os.Exit(1)
}
