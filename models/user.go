package models

import (
	"errors"
	"os"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Admin    bool   `json:"admin"`
}

// AuthenticateUser checks the provided username and password against the file
func AuthenticateUser(username, password string) (*User, error) {
	envUsername := os.Getenv("USERNAME")
	envPassword := os.Getenv("PASSWORD")
	if username != envUsername || password != envPassword {
		return nil, ErrInvalidCredentials
	}

	return &User{
		Username: username,
		Admin: true,
	}, nil
}