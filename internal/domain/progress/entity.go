package progress

import (
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("learning progress not found")
	ErrInvalidInput     = errors.New("invalid learning progress input")
	ErrDocumentNotFound = errors.New("document not found")
)

type LearningProgress struct {
	DocumentSlug      string          `json:"documentSlug"`
	Status            string          `json:"status"`
	ScrollPosition    int             `json:"scrollPosition"`
	Note              string          `json:"note"`
	CheckedFlashcards map[string]bool `json:"checkedFlashcards"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}
