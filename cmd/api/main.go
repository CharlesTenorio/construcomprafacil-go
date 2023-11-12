package main

import (
	"log"
	"net/http"

	"github.com/katana/back-end/construcomprafacil-go/internal/config"
	"github.com/katana/back-end/construcomprafacil-go/internal/config/logger"

	"github.com/katana/back-end/construcomprafacil-go/pkg/server"

	"github.com/go-chi/chi/v5"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {

	logger.Info("start Hoodid application")
	conf := config.NewConfig()

	r := chi.NewRouter()

	r.Get("/", healthcheck)

	srv := server.NewHTTPServer(r, conf)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server Run on [Port: %s], [Mode: %s], [Version: %s], [Commit: %s]", conf.PORT, conf.Mode, VERSION, COMMIT)

	select {}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"MSG": "Server Ok", "codigo": 200}`))
}
