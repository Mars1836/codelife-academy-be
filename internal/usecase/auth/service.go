package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strconv"
	"strings"
	"time"

	domain "codelife-study-be/internal/domain/auth"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository domain.Repository
	mailer     domain.Mailer
	secret     []byte
	otpTTL     time.Duration
	tokenTTL   time.Duration
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(repository domain.Repository, mailer domain.Mailer, secret string, otpTTL, tokenTTL time.Duration) *Service {
	if secret == "" {
		secret = "dev-only-change-me"
	}
	if otpTTL <= 0 {
		otpTTL = 10 * time.Minute
	}
	if tokenTTL <= 0 {
		tokenTTL = 24 * time.Hour
	}
	return &Service{repository: repository, mailer: mailer, secret: []byte(secret), otpTTL: otpTTL, tokenTTL: tokenTTL}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) error {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return err
	}
	if len(input.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := s.repository.CreateUser(ctx, domain.User{ID: randomID(), Email: email, PasswordHash: string(hash)})
	if errors.Is(err, domain.ErrEmailAlreadyExists) {
		existing, findErr := s.repository.FindUserByEmail(ctx, email)
		if findErr != nil {
			return findErr
		}
		if existing.EmailVerified {
			return domain.ErrEmailAlreadyExists
		}
		return s.sendOTP(ctx, existing)
	}
	if err != nil {
		return err
	}
	return s.sendOTP(ctx, user)
}

func (s *Service) VerifyEmail(ctx context.Context, input VerifyEmailInput) (domain.User, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return domain.User{}, err
	}
	otp := strings.TrimSpace(input.OTP)
	if len(otp) != 6 {
		return domain.User{}, domain.ErrInvalidOTP
	}
	otpHash := s.hashOTP(email, otp)
	userID, err := s.repository.FindValidEmailOTP(ctx, email, otpHash, time.Now())
	if err != nil {
		return domain.User{}, domain.ErrInvalidOTP
	}
	if err := s.repository.MarkEmailVerified(ctx, userID); err != nil {
		return domain.User{}, err
	}
	if err := s.repository.ConsumeEmailOTP(ctx, userID, otpHash); err != nil {
		return domain.User{}, err
	}
	return s.repository.FindUserByID(ctx, userID)
}

func (s *Service) Login(ctx context.Context, input LoginInput) (domain.Session, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return domain.Session{}, domain.ErrInvalidCredentials
	}
	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return domain.Session{}, domain.ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		return domain.Session{}, domain.ErrInvalidCredentials
	}
	if !user.EmailVerified {
		return domain.Session{}, domain.ErrEmailNotVerified
	}
	token, err := s.signToken(user)
	if err != nil {
		return domain.Session{}, err
	}
	return domain.Session{User: user, Token: token}, nil
}

func (s *Service) Me(ctx context.Context, token string) (domain.User, error) {
	claims, err := s.verifyToken(token)
	if err != nil {
		return domain.User{}, domain.ErrInvalidCredentials
	}
	return s.repository.FindUserByID(ctx, claims.Subject)
}

func (s *Service) sendOTP(ctx context.Context, user domain.User) error {
	otp, err := generateOTP()
	if err != nil {
		return err
	}
	if err := s.repository.SaveEmailOTP(ctx, user.ID, user.Email, s.hashOTP(user.Email, otp), time.Now().Add(s.otpTTL)); err != nil {
		return err
	}
	if s.mailer == nil {
		return nil
	}
	return s.mailer.SendOTP(ctx, user.Email, otp)
}

func (s *Service) hashOTP(email, otp string) string {
	sum := sha256.Sum256([]byte(email + ":" + otp + ":" + string(s.secret)))
	return hex.EncodeToString(sum[:])
}

func normalizeEmail(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "", fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return "", fmt.Errorf("email is invalid")
	}
	return value, nil
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func randomID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(raw)
}

type tokenClaims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Expires int64  `json:"exp"`
}

func (s *Service) signToken(user domain.User) (string, error) {
	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(tokenClaims{Subject: user.ID, Email: user.Email, Expires: time.Now().Add(s.tokenTTL).Unix()})
	if err != nil {
		return "", err
	}
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	return unsigned + "." + s.signature(unsigned), nil
}

func (s *Service) verifyToken(token string) (tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return tokenClaims{}, errors.New("invalid token")
	}
	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(s.signature(unsigned))) {
		return tokenClaims{}, errors.New("invalid signature")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return tokenClaims{}, err
	}
	var claims tokenClaims
	if err := json.Unmarshal(raw, &claims); err != nil {
		return tokenClaims{}, err
	}
	if claims.Subject == "" || claims.Expires < time.Now().Unix() {
		return tokenClaims{}, errors.New("expired token")
	}
	return claims, nil
}

func (s *Service) signature(unsigned string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
