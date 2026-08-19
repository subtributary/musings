package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func addCMSRoutes(router *chi.Mux, deps *Dependencies, cfg Config) {
	router.Route("/_cms", func(r chi.Router) {
		r.Get("/files", tbdHandler())
		r.Get("/index", tbdHandler())
		r.Get("/config", tbdHandler())
	})
}

func tbdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//
	}
}
