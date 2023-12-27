package categoria

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/categoria"
)

func RegisterCategoriaPIHandlers(r chi.Router, service categoria.CategoriaServiceInterface) {
	r.Route("/api/v1/categoria", func(r chi.Router) {
		r.Post("/add", createCategoria(service))
		r.Put("/update/{id}/{nome}", updateCategoria(service))
		r.Get("/getbyid/{id}", getByIdCategoria(service))
		r.Get("/listprd/{id}", getProdutosPorCategoria(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllCategoria(service)
			handler.ServeHTTP(w, r)
		})
	})
}
