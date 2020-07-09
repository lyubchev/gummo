package main

import (
	"github.com/go-chi/chi"
)

type Web struct {
	Router *chi.Mux
}

func NewWeb() *Web {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", UserRegister)
		// r.Post("/login")
	})

	return &Web{Router: r}
}
