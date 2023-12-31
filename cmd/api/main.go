package main

import (
	"log"
	"net/http"

	"github.com/katana/back-end/orcafacil-go/internal/config"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	hand_categoria "github.com/katana/back-end/orcafacil-go/internal/handler/categoria"
	hand_meiopg "github.com/katana/back-end/orcafacil-go/internal/handler/meiospg"
	hand_prd "github.com/katana/back-end/orcafacil-go/internal/handler/produto"
	hand_usr "github.com/katana/back-end/orcafacil-go/internal/handler/user"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"

	"github.com/katana/back-end/orcafacil-go/pkg/server"
	service_categoria "github.com/katana/back-end/orcafacil-go/pkg/service/categoria"
	service_meiopg "github.com/katana/back-end/orcafacil-go/pkg/service/meiospg"
	service_prd "github.com/katana/back-end/orcafacil-go/pkg/service/produto"
	service_usr "github.com/katana/back-end/orcafacil-go/pkg/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	usr_service := service_usr.NewUsuarioservice(mogDbConn)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", conf.TokenAuth))
	r.Use(middleware.WithValue("JWTTokenExp", conf.JWTTokenExp))

	// Adicione o middleware JwtMiddleware para autenticação JWT

	r.Get("/", healthcheck)
	hand_meiopg.RegisterMeioPgAPIHandlers(r, meiopg_service)
	hand_categoria.RegisterCategoriaPIHandlers(r, categoria_service)
	hand_prd.RegisterProdutoAPIHandlers(r, prd_service)
	hand_usr.RegisterUsuarioAPIHandlers(r, usr_service)

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
