package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var ErrBadRequest = errors.New("Bad Request")

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var registerForm registerRequest

	if err := render.DecodeJSON(r.Body, registerForm); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}
}
