package auth

import (
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidOTP         = errors.New("invalid otp")
	ErrNotFound           = errors.New("auth record not found")
)

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Session struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
