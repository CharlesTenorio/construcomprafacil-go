package orcamento

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/orcamento"
)

func RegisterOrcamentoPIHandlers(r chi.Router, service orcamento.OrcamentoServiceInterface) {
	r.Route("/api/v1/cliente", func(r chi.Router) {
		r.Post("/add", createorcamento(service))
		r.Put("/update/{id}/{nome}", updateorcamento(service))
		r.Get("/getbyid/{id}", getByIdorcamento(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllorcamento(service)
			handler.ServeHTTP(w, r)
		})
	})
}
