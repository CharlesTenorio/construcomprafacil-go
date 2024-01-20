package catalago

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/catalagoprd"
)

func RegisterCatalogoPIHandlers(r chi.Router, service catalagoprd.CatalagoServiceInterface) {
	r.Route("/api/v1/catalogo", func(r chi.Router) {

		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllCatalago(service)
			handler.ServeHTTP(w, r)
		})
	})
}
