package progress

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domain "codelife-study-be/internal/domain/progress"
)

type Service struct {
	repository domain.Repository
}

type UpdateInput struct {
	Status            *string          `json:"status"`
	ScrollPosition    *int             `json:"scrollPosition"`
	Note              *string          `json:"note"`
	CheckedFlashcards *map[string]bool `json:"checkedFlashcards"`
}

func New(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, userID string) ([]domain.LearningProgress, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("%w: user id is required", domain.ErrInvalidInput)
	}
	return s.repository.ListByUser(ctx, userID)
}

func (s *Service) Update(ctx context.Context, userID, documentSlug string, input UpdateInput) (domain.LearningProgress, error) {
	documentSlug = strings.TrimSpace(documentSlug)
	if userID == "" || documentSlug == "" || strings.ContainsAny(documentSlug, `/\\`) {
		return domain.LearningProgress{}, fmt.Errorf("%w: invalid target", domain.ErrInvalidInput)
	}

	progress, err := s.repository.FindByUserAndDocument(ctx, userID, documentSlug)
	if errors.Is(err, domain.ErrNotFound) {
		progress = domain.LearningProgress{
			DocumentSlug:      documentSlug,
			Status:            "unread",
			CheckedFlashcards: map[string]bool{},
		}
	} else if err != nil {
		return domain.LearningProgress{}, err
	}

	if input.Status != nil {
		if !validStatus(*input.Status) {
			return domain.LearningProgress{}, fmt.Errorf("%w: status must be unread, studying, or completed", domain.ErrInvalidInput)
		}
		progress.Status = *input.Status
	}
	if input.ScrollPosition != nil {
		if *input.ScrollPosition < 0 {
			return domain.LearningProgress{}, fmt.Errorf("%w: scroll position cannot be negative", domain.ErrInvalidInput)
		}
		progress.ScrollPosition = *input.ScrollPosition
	}
	if input.Note != nil {
		if len(*input.Note) > 50000 {
			return domain.LearningProgress{}, fmt.Errorf("%w: note is too long", domain.ErrInvalidInput)
		}
		progress.Note = *input.Note
	}
	if input.CheckedFlashcards != nil {
		if len(*input.CheckedFlashcards) > 100 {
			return domain.LearningProgress{}, fmt.Errorf("%w: too many flashcards", domain.ErrInvalidInput)
		}
		progress.CheckedFlashcards = *input.CheckedFlashcards
	}

	return s.repository.Upsert(ctx, userID, progress)
}

func validStatus(status string) bool {
	return status == "unread" || status == "studying" || status == "completed"
}
