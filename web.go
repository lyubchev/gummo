package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-redis/redis"
	"upper.io/db.v3"
)

const (
	UserContextKey = "gummo_user"
)

var (
	ErrBadTypeAssertion = errors.New("failed to assert type")
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
		r.With(wb.Authenticator).Get("/me", wb.UserAbout)
		r.With(wb.Authenticator).Get("/logout", wb.UserLogout)
	})

	return wb
}

func (wb *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wb.Router.ServeHTTP(w, r)
}

func (wb *Web) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := wb.GetFromRequest(r)

		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, http.StatusText(http.StatusUnauthorized))
			return
		}

		ctx := WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Get gets a session from Redis based on the session ID
func (wb *Web) Get(sID string) (*User, error) {
	email, err := wb.Redis.Get(sID).Result()
	if err != nil {
		return &User{}, err
	}

	var user User
	err = wb.DB.Collection("users").Find("email", email).One(&user)
	return &user, err
}

// GetFromRequest gets a session from Redis based on the Cookie value from the request
func (wb *Web) GetFromRequest(r *http.Request) (*User, error) {
	cookie, err := r.Cookie(CookieKey)
	if err != nil {
		return nil, err
	}

	return wb.Get(cookie.Value)
}

// WithUser sets a user to a context
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserContextKey, *user)
}

// CurrentUser gets a user from a context
func CurrentUser(ctx context.Context) (User, error) {
	if user, ok := ctx.Value(UserContextKey).(User); ok {
		return user, nil
	}

	return User{}, ErrBadTypeAssertion
}
