package document

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("document not found")

// Repository is owned by the domain. Adapters implement it.
type Repository interface {
	List(context.Context) ([]Document, error)
	FindBySlug(context.Context, string) (Document, error)
}

type Cache interface {
	Get(context.Context, string, any) bool
	Set(context.Context, string, any) error
}
