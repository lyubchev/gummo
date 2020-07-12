package main

import (
	"github.com/go-chi/chi"
	"upper.io/db.v3"
)

type Web struct {
	Router *chi.Mux
	DB     db.Database
}

func NewWeb(db db.Database) *Web {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", UserRegister)
		// r.Post("/login")
	})

	return &Web{Router: r, DB: db}
}
