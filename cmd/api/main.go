package main

import (
	"log"
	"net/http"

	"github.com/katana/back-end/orcafacil-go/internal/config"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	hand_categoria "github.com/katana/back-end/orcafacil-go/internal/handler/categoria"
	hand_meiopg "github.com/katana/back-end/orcafacil-go/internal/handler/meiospg"
	hand_prd "github.com/katana/back-end/orcafacil-go/internal/handler/produto"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"

	"github.com/katana/back-end/orcafacil-go/pkg/server"
	service_categoria "github.com/katana/back-end/orcafacil-go/pkg/service/categoria"
	service_meiopg "github.com/katana/back-end/orcafacil-go/pkg/service/meiospg"
	service_prd "github.com/katana/back-end/orcafacil-go/pkg/service/produto"

	"github.com/go-chi/chi/v5"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {

	logger.Info("start Application Cota Facil")
	conf := config.NewConfig()

	mogDbConn := mongodb.New(conf)
	meiopg_service := service_meiopg.NewMeioPgService(mogDbConn)
	categoria_service := service_categoria.NewCategoriaervice(mogDbConn)
	prd_service := service_prd.NewProdutoervice(mogDbConn)

	r := chi.NewRouter()

	r.Get("/", healthcheck)
	hand_meiopg.RegisterMeioPgAPIHandlers(r, meiopg_service)
	hand_categoria.RegisterCategoriaPIHandlers(r, categoria_service)
	hand_prd.RegisterProdutoAPIHandlers(r, prd_service)

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
