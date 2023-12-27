package produto

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/produto"
)

func RegisterProdutoAPIHandlers(r chi.Router, service produto.ProdutoServiceInterface) {
	r.Route("/api/v1/produto", func(r chi.Router) {
		r.Post("/add", createProduto(service))
		r.Put("/update/{id}/{nome}", updateProduto(service))
		r.Get("/getbyid/{id}", getByIdProduto(service)) // Adicionado a barra no início
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllProduto(service)
			handler.ServeHTTP(w, r)
		})
	})
}
