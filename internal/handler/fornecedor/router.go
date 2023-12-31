package fornecedor

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/fornecedor"
)

func RegisterFornecedorPIHandlers(r chi.Router, service fornecedor.FornecedorServiceInterface) {
	r.Route("/api/v1/fornecedor", func(r chi.Router) {
		r.Post("/add", createFornecedor(service))
		r.Put("/update/{id}/{nome}", updateFornecedor(service))
		r.Get("/getbyid/{id}", getByIdFornecedor(service))
		r.Get("/listprd/{id}", getProdutosPorFornecedor(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllFornecedor(service)
			handler.ServeHTTP(w, r)
		})
	})
}
