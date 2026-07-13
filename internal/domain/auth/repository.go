package auth

import (
	"context"
	"time"
)

type Repository interface {
	CreateUser(context.Context, User) (User, error)
	FindUserByEmail(context.Context, string) (User, error)
	FindUserByID(context.Context, string) (User, error)
	MarkEmailVerified(context.Context, string) error
	SaveEmailOTP(ctx context.Context, userID, email, otpHash string, expiresAt time.Time) error
	FindValidEmailOTP(ctx context.Context, email, otpHash string, now time.Time) (string, error)
	ConsumeEmailOTP(ctx context.Context, userID, otpHash string) error
}

type Mailer interface {
	SendOTP(context.Context, string, string) error
}
