package subcategoria

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/subcategoria"
)

func RegisterSbuCategoriaPIHandlers(r chi.Router, service subcategoria.SubcategoriaServiceInterface) {
	r.Route("/api/v1/subcategoria", func(r chi.Router) {
		r.Post("/add", createSubCategoria(service))
		r.Put("/update/{id}", updateSubCategoria(service))
		r.Get("/getbyid/{id}", getByIdSubCategoria(service))

		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllSubCategoria(service)
			handler.ServeHTTP(w, r)
		})
	})
}
