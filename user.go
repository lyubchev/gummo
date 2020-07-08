package main

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string
	Password string
	Name     string
	Avatar   string
	ID       string
	Location Location
}

func New(email, password, name, avatar string) (*User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:    email,
		Password: hashed,
		Name:     name,
		Avatar:   avatar,
		ID:       uuid.New().String(),
	}, nil
}

func hashPassword(pwd string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPwd), nil
}
