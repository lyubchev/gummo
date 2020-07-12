package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var ErrBadRequest = errors.New("Bad Request")
var ErrNotFound = errors.New("Not found")
var ErrInternalServer = errors.New("Internal server error")

type regCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
}

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (web *Web) UserRegister(w http.ResponseWriter, r *http.Request) {
	var creds regCredentials

	if err := render.DecodeJSON(r.Body, creds); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrInternalServer)
		return
	}

	user, err := NewUser(creds.Email, creds.Password, creds.Name, creds.Avatar)
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

func (wb *Web) UserLogin(w http.ResponseWriter, r *http.Request) {
	var creds loginCredentials

	if err := render.DecodeJSON(r.Body, creds); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrInternalServer)
		return
	}

	var user User

	err := wb.DB.
		Collection("users").
		Find("email", creds.Email).One(&user)

	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, ErrNotFound)
		return
	}

	if user.CompareHashAndPassword(creds.Password) {
		sid := uuid.New().String()
		if err := wb.Redis.Set(sid, user.Email, time.Hour*24); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrInternalServer)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "gummo_token",
			Value:   sid,
			Expires: time.Now().Add(time.Hour * 24),
		})

		render.Status(r, http.StatusOK)
	}

}
