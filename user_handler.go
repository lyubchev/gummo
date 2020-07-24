package main

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

const (
	CookieKey = "gummo_token"
)

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

	if err := render.DecodeJSON(r.Body, &creds); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := NewUser(creds.Email, creds.Password, creds.Name, creds.Avatar)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, http.StatusText(http.StatusBadRequest))
		return
	}

	_, err = web.DB.Collection("users").Insert(user)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				render.Status(r, http.StatusUnprocessableEntity)
				render.JSON(w, r, http.StatusText(http.StatusUnprocessableEntity))
				return
			}
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, http.StatusText(http.StatusInternalServerError))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (wb *Web) UserLogin(w http.ResponseWriter, r *http.Request) {
	var creds loginCredentials

	if err := render.DecodeJSON(r.Body, &creds); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, http.StatusText(http.StatusInternalServerError))
		return
	}

	var user User

	err := wb.DB.
		Collection("users").
		Find("email", creds.Email).One(&user)

	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, http.StatusText(http.StatusNotFound))
		return
	}

	if !user.CompareHashAndPassword(creds.Password) {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, http.StatusText(http.StatusUnauthorized))
		return
	}

	sid := uuid.New().String()
	if err := wb.Redis.Set(sid, user.Email, time.Hour*24).Err(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    CookieKey,
		Value:   sid,
		Expires: time.Now().Add(time.Hour * 24),
	})

	render.Status(r, http.StatusOK)

}

func (wb *Web) UserAbout(w http.ResponseWriter, r *http.Request) {
	user, err := CurrentUser(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, http.StatusText(http.StatusInternalServerError))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}
