package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var ErrBadRequest = errors.New("Bad Request")
var ErrInternalServer = errors.New("Internal server error")

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
}

func (web *Web) UserRegister(w http.ResponseWriter, r *http.Request) {
	var registerForm registerRequest

	if err := render.DecodeJSON(r.Body, registerForm); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrInternalServer)
		return
	}

	user, err := NewUser(registerForm.Email, registerForm.Password, registerForm.Name, registerForm.Avatar)
	if err != nil {

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

	_, err = web.DB.Collection("users").Insert(user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}
