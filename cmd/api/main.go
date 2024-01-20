package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	hand_catalago "github.com/katana/back-end/orcafacil-go/internal/handler/catalago"
	hand_categoria "github.com/katana/back-end/orcafacil-go/internal/handler/categoria"
	hand_cliente "github.com/katana/back-end/orcafacil-go/internal/handler/cliente"
	hand_fornec "github.com/katana/back-end/orcafacil-go/internal/handler/fornecedor"
	hand_meiopg "github.com/katana/back-end/orcafacil-go/internal/handler/meiospg"
	hand_orca "github.com/katana/back-end/orcafacil-go/internal/handler/orcamento"
	hand_produto "github.com/katana/back-end/orcafacil-go/internal/handler/produto"
	hand_usr "github.com/katana/back-end/orcafacil-go/internal/handler/user"

	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/rabbitmq"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/redisdb"

	"github.com/katana/back-end/orcafacil-go/pkg/server"
	"github.com/katana/back-end/orcafacil-go/pkg/service/catalagoprd"
	service_catalago "github.com/katana/back-end/orcafacil-go/pkg/service/catalagoprd"
	service_usr "github.com/katana/back-end/orcafacil-go/pkg/service/user"

	service_categoria "github.com/katana/back-end/orcafacil-go/pkg/service/categoria"
	service_cliente "github.com/katana/back-end/orcafacil-go/pkg/service/cliente"
	service_produto "github.com/katana/back-end/orcafacil-go/pkg/service/produto"

	service_fornec "github.com/katana/back-end/orcafacil-go/pkg/service/fornecedor"
	service_meiopg "github.com/katana/back-end/orcafacil-go/pkg/service/meiospg"
	service_orcamento "github.com/katana/back-end/orcafacil-go/pkg/service/orcamento"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {
	fila := []rabbitmq.Fila{
		{
			Name:    "QUEUE_PRDS_PARA_COTACAO",
			Durable: true,
		},
	}
	logger.Info("start Application Cota Facil")
	conf := config.NewConfig()

	mogDbConn := mongodb.New(conf)
	rbtMQConn := rabbitmq.NewRabbitMQ(fila, conf)
	rdisConn := redisdb.NewRedisClient(conf)
	usr_service := service_usr.NewUsuarioservice(mogDbConn)
	meiopg_service := service_meiopg.NewMeioPgService(mogDbConn)
	categoria_service := service_categoria.NewCategoriaervice(mogDbConn)
	prd_service := service_produto.NewProdutoervice(mogDbConn)

	fornec_service := service_fornec.NewFornecedorervice(mogDbConn)
	cli_service := service_cliente.NewClienteervice(mogDbConn)
	orca_service := service_orcamento.NewOrcamentoService(mogDbConn, rbtMQConn)
	catalago_service := service_catalago.NewCatalagoService(mogDbConn, rdisConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", conf.TokenAuth))
	r.Use(middleware.WithValue("JWTTokenExp", conf.JWTTokenExp))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", healthcheck)
	hand_usr.RegisterUsuarioAPIHandlers(r, usr_service)
	hand_meiopg.RegisterMeioPgAPIHandlers(r, meiopg_service, conf)
	hand_produto.RegisterPrdPIHandlers(r, prd_service)
	hand_categoria.RegisterCategoriaPIHandlers(r, categoria_service)

	hand_fornec.RegisterFornecedorPIHandlers(r, fornec_service)
	hand_cliente.RegisterClientePIHandlers(r, cli_service)
	hand_orca.RegisterOrcamentoPIHandlers(r, orca_service)
	hand_catalago.RegisterCatalogoPIHandlers(r, catalago_service)

	catalogService := catalagoprd.NewCatalagoService(mogDbConn, rdisConn)

	// Inicie o worker em uma goroutine
	go startWorker(catalogService)

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

func startWorker(catalogService catalagoprd.CatalagoServiceInterface) {
	// Crie um canal para sinalizar quando o worker deve parar
	stopCh := make(chan struct{})

	// Inicie o worker em uma goroutine
	go func() {
		// Loop infinito para executar a função SalveRedis a cada 5 minutos
		for {
			select {
			case <-stopCh:
				// Sinal recebido para parar o worker
				return
			default:
				// Execute a função SalveRedis
				_, err := catalogService.SalveRedis(context.Background())
				if err != nil {
					logger.Error("Erro ao executar SalveRedis:", err)
				} else {
					logger.Info("SalveRedis executado com sucesso")
				}

				// Aguarde 5 minutos antes de executar novamente
				time.Sleep(5 * time.Minute)
			}
		}
	}()

	// Aguarde um sinal para parar o worker (pode ser um sinal do sistema operacional, por exemplo)
	<-stopCh
	log.Println("Worker parado")
}
