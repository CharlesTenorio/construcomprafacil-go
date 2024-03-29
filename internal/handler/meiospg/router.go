package meiospg

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/katana/back-end/orcafacil-go/internal/config"
	"github.com/katana/back-end/orcafacil-go/pkg/service/meiospg"
)

func RegisterMeioPgAPIHandlers(r chi.Router, service meiospg.MeiosServiceInterface, cf *config.Config) {

	r.Route("/api/v1/meiopag", func(r chi.Router) {
		r.Use(jwtauth.Verifier(cf.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/add", createMeioPg(service))
		r.Put("/update/{id}/{nome}", updateMeioPg(service))
		r.Get("/getbyid/{id}", getByIdMeioPg(service)) // Adicionado a barra no início
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllMeioPg(service)
			handler.ServeHTTP(w, r)
		})
	})
}
