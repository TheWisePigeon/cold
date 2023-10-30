package handlers

import (
	"cold/handlers/auth"
	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(r *chi.Mux) {
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", auth.Register)
	})

  r.Get("/register", auth.GetRegistrationPage)
}
