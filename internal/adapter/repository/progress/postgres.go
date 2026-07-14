package progress

import (
	"context"
	"encoding/json"
	"errors"

	domain "codelife-study-be/internal/domain/progress"

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

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]domain.LearningProgress, error) {
	rows, err := r.db.Query(ctx, `
		SELECT document_slug, status, scroll_position, note, checked_flashcards, updated_at
		FROM user_document_progress
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, scanProgress)
}

func (r *PostgresRepository) FindByUserAndDocument(ctx context.Context, userID, documentSlug string) (domain.LearningProgress, error) {
	progress, err := scanProgress(r.db.QueryRow(ctx, `
		SELECT document_slug, status, scroll_position, note, checked_flashcards, updated_at
		FROM user_document_progress
		WHERE user_id = $1 AND document_slug = $2
	`, userID, documentSlug))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.LearningProgress{}, domain.ErrNotFound
	}
	return progress, err
}

func (r *PostgresRepository) Upsert(ctx context.Context, userID string, progress domain.LearningProgress) (domain.LearningProgress, error) {
	checked, err := json.Marshal(progress.CheckedFlashcards)
	if err != nil {
		return domain.LearningProgress{}, err
	}
	result, err := scanProgress(r.db.QueryRow(ctx, `
		INSERT INTO user_document_progress (
			user_id, document_slug, status, scroll_position, note, checked_flashcards, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (user_id, document_slug) DO UPDATE SET
			status = EXCLUDED.status,
			scroll_position = EXCLUDED.scroll_position,
			note = EXCLUDED.note,
			checked_flashcards = EXCLUDED.checked_flashcards,
			updated_at = NOW()
		RETURNING document_slug, status, scroll_position, note, checked_flashcards, updated_at
	`, userID, progress.DocumentSlug, progress.Status, progress.ScrollPosition, progress.Note, checked))
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return domain.LearningProgress{}, domain.ErrDocumentNotFound
	}
	return result, err
}

func scanProgress(row pgx.CollectableRow) (domain.LearningProgress, error) {
	var progress domain.LearningProgress
	var checked []byte
	if err := row.Scan(
		&progress.DocumentSlug,
		&progress.Status,
		&progress.ScrollPosition,
		&progress.Note,
		&checked,
		&progress.UpdatedAt,
	); err != nil {
		return domain.LearningProgress{}, err
	}
	if err := json.Unmarshal(checked, &progress.CheckedFlashcards); err != nil {
		return domain.LearningProgress{}, err
	}
	if progress.CheckedFlashcards == nil {
		progress.CheckedFlashcards = map[string]bool{}
	}
	return progress, nil
}
