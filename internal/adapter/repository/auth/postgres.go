package auth

import (
	"context"
	"errors"
	"time"

	domain "codelife-study-be/internal/domain/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO auth_users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, email_verified, created_at
	`, user.ID, user.Email, user.PasswordHash).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified, &user.CreatedAt)
	if isUniqueViolation(err) {
		return domain.User{}, domain.ErrEmailAlreadyExists
	}
	return user, err
}

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.findUser(ctx, "email", email)
}

func (r *PostgresRepository) FindUserByID(ctx context.Context, id string) (domain.User, error) {
	return r.findUser(ctx, "id", id)
}

func (r *PostgresRepository) MarkEmailVerified(ctx context.Context, userID string) error {
	tag, err := r.db.Exec(ctx, `UPDATE auth_users SET email_verified = TRUE, updated_at = NOW() WHERE id = $1`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) SaveEmailOTP(ctx context.Context, userID, email, otpHash string, expiresAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE auth_email_otps SET consumed_at = NOW() WHERE user_id = $1 AND consumed_at IS NULL`, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO auth_email_otps (user_id, email, otp_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, userID, email, otpHash, expiresAt); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) FindValidEmailOTP(ctx context.Context, email, otpHash string, now time.Time) (string, error) {
	var userID string
	err := r.db.QueryRow(ctx, `
		SELECT user_id
		FROM auth_email_otps
		WHERE email = $1 AND otp_hash = $2 AND consumed_at IS NULL AND expires_at > $3
		ORDER BY created_at DESC
		LIMIT 1
	`, email, otpHash, now).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrInvalidOTP
	}
	return userID, err
}

func (r *PostgresRepository) ConsumeEmailOTP(ctx context.Context, userID, otpHash string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE auth_email_otps
		SET consumed_at = NOW()
		WHERE user_id = $1 AND otp_hash = $2 AND consumed_at IS NULL
	`, userID, otpHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrInvalidOTP
	}
	return nil
}

func (r *PostgresRepository) findUser(ctx context.Context, column, value string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, password_hash, email_verified, created_at
		FROM auth_users
		WHERE ` + column + ` = $1
	`
	err := r.db.QueryRow(ctx, query, value).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	return user, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
