package auth

import (
	"context"
	"testing"
	"time"

	domain "codelife-study-be/internal/domain/auth"
)

type fakeRepository struct {
	users map[string]domain.User
	otps  map[string]otpRecord
}

type otpRecord struct {
	userID    string
	otpHash   string
	expiresAt time.Time
	consumed  bool
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{users: map[string]domain.User{}, otps: map[string]otpRecord{}}
}

func (r *fakeRepository) CreateUser(_ context.Context, user domain.User) (domain.User, error) {
	if _, ok := r.users[user.Email]; ok {
		return domain.User{}, domain.ErrEmailAlreadyExists
	}
	user.CreatedAt = time.Now()
	r.users[user.Email] = user
	return user, nil
}

func (r *fakeRepository) FindUserByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := r.users[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (r *fakeRepository) FindUserByID(_ context.Context, id string) (domain.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (r *fakeRepository) MarkEmailVerified(_ context.Context, userID string) error {
	for email, user := range r.users {
		if user.ID == userID {
			user.EmailVerified = true
			r.users[email] = user
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *fakeRepository) SaveEmailOTP(_ context.Context, userID, email, otpHash string, expiresAt time.Time) error {
	r.otps[email] = otpRecord{userID: userID, otpHash: otpHash, expiresAt: expiresAt}
	return nil
}

func (r *fakeRepository) FindValidEmailOTP(_ context.Context, email, otpHash string, now time.Time) (string, error) {
	record, ok := r.otps[email]
	if !ok || record.consumed || record.otpHash != otpHash || !record.expiresAt.After(now) {
		return "", domain.ErrInvalidOTP
	}
	return record.userID, nil
}

func (r *fakeRepository) ConsumeEmailOTP(_ context.Context, userID, otpHash string) error {
	for email, record := range r.otps {
		if record.userID == userID && record.otpHash == otpHash {
			record.consumed = true
			r.otps[email] = record
			return nil
		}
	}
	return domain.ErrInvalidOTP
}

type fakeMailer struct {
	email string
	otp   string
}

func (m *fakeMailer) SendOTP(_ context.Context, email, otp string) error {
	m.email = email
	m.otp = otp
	return nil
}

func TestRegisterVerifyAndLogin(t *testing.T) {
	repository := newFakeRepository()
	mailer := &fakeMailer{}
	service := New(repository, mailer, "secret", time.Minute, time.Hour)

	if err := service.Register(context.Background(), RegisterInput{Email: "USER@example.com", Password: "password123"}); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if mailer.email != "user@example.com" || len(mailer.otp) != 6 {
		t.Fatalf("unexpected otp mail: %#v", mailer)
	}
	if _, err := service.Login(context.Background(), LoginInput{Email: "user@example.com", Password: "password123"}); err != domain.ErrEmailNotVerified {
		t.Fatalf("expected email not verified, got %v", err)
	}
	user, err := service.VerifyEmail(context.Background(), VerifyEmailInput{Email: "user@example.com", OTP: mailer.otp})
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if !user.EmailVerified {
		t.Fatalf("expected verified user")
	}
	session, err := service.Login(context.Background(), LoginInput{Email: "user@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if session.Token == "" || session.User.Email != "user@example.com" {
		t.Fatalf("unexpected session: %#v", session)
	}
}

func TestVerifyRejectsInvalidOTP(t *testing.T) {
	repository := newFakeRepository()
	service := New(repository, &fakeMailer{}, "secret", time.Minute, time.Hour)
	if err := service.Register(context.Background(), RegisterInput{Email: "user@example.com", Password: "password123"}); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if _, err := service.VerifyEmail(context.Background(), VerifyEmailInput{Email: "user@example.com", OTP: "000000"}); err != domain.ErrInvalidOTP {
		t.Fatalf("expected invalid otp, got %v", err)
	}
}

func TestRegisterAgainResendsOTPForUnverifiedEmail(t *testing.T) {
	repository := newFakeRepository()
	mailer := &fakeMailer{}
	service := New(repository, mailer, "secret", time.Minute, time.Hour)
	if err := service.Register(context.Background(), RegisterInput{Email: "user@example.com", Password: "password123"}); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	firstOTP := mailer.otp
	if err := service.Register(context.Background(), RegisterInput{Email: "user@example.com", Password: "password123"}); err != nil {
		t.Fatalf("register resend failed: %v", err)
	}
	if mailer.otp == "" || firstOTP == "" {
		t.Fatalf("expected otp values")
	}
}
