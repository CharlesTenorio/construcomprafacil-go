package meiospg

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/meiospg"
)

func RegisterMeioPgAPIHandlers(r chi.Router, service meiospg.MeiosServiceInterface) {
	r.Route("/api/v1/meiopg", func(r chi.Router) {
		r.Post("/add", createMeioPg(service))
		r.Put("/update/{id}/{nome}", updateMeioPg(service))
		r.Get("getbyid/{id}", getByIdMeioPg(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllMeioPg(service)
			handler.ServeHTTP(w, r)
		})
	})

}
