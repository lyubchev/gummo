package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"upper.io/db.v3"
)

type Web struct {
	Router *chi.Mux
	DB     db.Database
	Redis  *redis.Client
}

func NewWeb(db db.Database, rdb *redis.Client) *Web {
	r := chi.NewRouter()
	wb := &Web{Router: r, DB: db, Redis: rdb}

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/ping"),
		middleware.RedirectSlashes,
		middleware.StripSlashes,
		middleware.Timeout(time.Minute),
	)

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", wb.UserRegister)
		r.Post("/login", wb.UserLogin)
		r.Get("/welcome", wb.UserWelcome)
	})

	return wb
}

func (wb *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wb.Router.ServeHTTP(w, r)
}
