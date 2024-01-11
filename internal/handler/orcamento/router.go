package orcamento

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/orcamento"
)

func RegisterOrcamentoPIHandlers(r chi.Router, service orcamento.OrcamentoServiceInterface) {
	r.Route("/api/v1/orcamento", func(r chi.Router) {
		r.Post("/add", createOrcamento(service))
		r.Put("/update/{id}/{nome}", updateOrcamento(service))
		r.Get("/getbyid/{id}", getByIdOrcamento(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllOrcamento(service)
			handler.ServeHTTP(w, r)
		})
	})
}
